package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const (
	// The sts GetCallerIdentity request is valid for 15 minutes regardless of this parameters value after it has been
	// signed, but we set this unused parameter to 60 for legacy reasons (we check for a value between 0 and 60 on the
	// server side in 0.3.0 or earlier).  IT IS IGNORED.  If we can get STS to support x-amz-expires, then we should
	// set this parameter to the actual expiration, and make it configurable.
	requestPresignParam = "60"
	// The actual token expiration (presigned STS urls are valid for 15 minutes after timestamp in x-amz-date).
	presignedURLExpiration = 15 * time.Minute
	v1Prefix               = "k8s-aws-v1."
	clusterIDHeader        = "x-k8s-aws-id"
	// Format of the X-Amz-Date header used for expiration
	// https://golang.org/pkg/time/#pkg-constants
	dateHeaderFormat = "20060102T150405Z"
)

var (
	awsConfig *aws.Config
	version   string
)

func main() {

	app := cli.NewApp()

	app.Before = func(c *cli.Context) error {
		awsConfigErr := error(nil)

		awsConfig, awsConfigErr = loadAWSConfig("eu-west-1")

		if awsConfigErr != nil {
			return awsConfigErr
		}

		return awsConfigErr
	}

	app.Version = version
	app.Authors = []*cli.Author{
		{
			Name:  "NinjaOps by https://raftech.nl",
			Email: "hello@raftech.nl",
		},
	}
	app.Name = "qbconf"
	app.Usage = "Minimalistic Kubernetes kubeconfig file generator using AWS STS and EKS APIs"

	app.Commands = []*cli.Command{
		{
			Name:  "generate-gha",
			Usage: "Generate a kubeconfig file for an EKS cluster by assuming specified AWS IAM role using GHA OIDC",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "role-arn",
					Usage:    "ARN of the AWS IAM role to assume",
					EnvVars:  []string{"AWS_ROLE_ARN"},
					Value:    "",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "region",
					Usage:    "AWS region",
					EnvVars:  []string{"AWS_REGION"},
					Value:    "eu-west-1",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "role-session-name",
					Usage:    "Name of the AWS STS role session to create",
					EnvVars:  []string{"AWS_ROLE_SESSION_NAME"},
					Value:    "qbconf-session",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "eks-cluster-name",
					Usage:    "Name of the EKS cluster to generate a kubeconfig file for",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "output-file",
					Usage:    "Name of the file to write the generated kubeconfig to",
					Value:    "kubeconfig.yaml",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {

				client := resty.New()

				// Retrieve the token and URL from the environment of GitHub Actions
				tokenRequestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
				tokenRequestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

				resp, err := client.R().
					EnableTrace().
					SetAuthToken(tokenRequestToken).
					Get(fmt.Sprintf("%s&audience=sts.amazonaws.com", tokenRequestURL))

				if err != nil {
					return err
				}

				tokenValue := gjson.Get(resp.String(), "value").String()

				//return TokenWithRoleFromArn(, , , c.String("eks-cluster-name"), c.String("output-file"))
				awsConfig.Credentials = assumeRoleWithWebIdentity(c.String("role-arn"), c.String("role-session-name"), c.String("region"), tokenValue)

				result, _ := getAWSIdentity(*awsConfig)
				fmt.Println(*result.Arn)

				errKubeConfig := generateKubeConfig(c.String("region"), c.String("eks-cluster-name"), c.String("output-file"))
				if errKubeConfig != nil {
					return errKubeConfig
				}

				return err
			},
		},
		{
			Name:  "generate",
			Usage: "Generate a kubeconfig file for an EKS cluster by assuming specified AWS IAM role",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "role-arn",
					Usage:    "ARN of the AWS IAM role to assume",
					EnvVars:  []string{"AWS_ROLE_ARN"},
					Value:    "",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "region",
					Usage:    "AWS region",
					EnvVars:  []string{"AWS_REGION"},
					Value:    "eu-west-1",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "role-session-name",
					Usage:    "Name of the AWS STS role session to create",
					EnvVars:  []string{"AWS_ROLE_SESSION_NAME"},
					Value:    "qbconf-session",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "eks-cluster-name",
					Usage:    "Name of the EKS cluster to generate a kubeconfig file for",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "output-file",
					Usage:    "Name of the file to write the generated kubeconfig to",
					Value:    "kubeconfig.yaml",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {

				provider := assumeRoleByArn(c.String("role-arn"), c.String("role-session-name"), awsConfig)
				awsConfig.Credentials = provider

				result, _ := getAWSIdentity(*awsConfig)
				fmt.Println(*result.Arn)

				errKubeConfig := generateKubeConfig(c.String("region"), c.String("eks-cluster-name"), c.String("output-file"))
				if errKubeConfig != nil {
					return errKubeConfig
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadAWSConfig(region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Errorf("failed to load SDK config: %v", err)
		return nil, err
	}

	return &cfg, nil
}

func getAWSIdentity(cfg aws.Config) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.NewFromConfig(cfg)

	input := &sts.GetCallerIdentityInput{}
	result, err := svc.GetCallerIdentity(context.TODO(), input)
	if err != nil {
		var aerr awserr.Error
		if ok := errors.As(err, &aerr); ok {
			switch aerr.Code() {
			default:
				return nil, fmt.Errorf("aws error: %s", aerr.Error())
			}
		} else {
			return nil, fmt.Errorf("unexpected error: %s", err.Error())
		}
	}

	return result, nil
}

func assumeRoleByArn(roleArn, roleSessionName string, awsConfig *aws.Config) *stscreds.AssumeRoleProvider {

	// Create an STS client using the default config
	stsClient := sts.NewFromConfig(*awsConfig)

	// Create an AssumeRoleProvider that will assume the specified role
	roleProvider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = roleSessionName
	})

	return roleProvider
}

func assumeRoleWithWebIdentity(roleArn, roleSessionName, region, token string) *aws.CredentialsCache {
	// Set up AWS SDK config with your AWS region and profile name
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	// Create an STS client using the default config
	stsClient := sts.NewFromConfig(cfg)

	// // Set up the IAM role ARN that you want to assume
	// roleArn, err = arn.Parse(roleArn)
	// if err != nil {
	// 	panic(err)
	// }

	// Set up the AssumeRoleWithWebIdentity input
	input := &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String(roleSessionName),
		WebIdentityToken: aws.String(token),
	}

	// Call the AssumeRoleWithWebIdentity API to assume the IAM role
	resp, err := stsClient.AssumeRoleWithWebIdentity(context.Background(), input)
	if err != nil {
		panic(err)
	}

	// value := aws.Credentials{
	// 	AccessKeyID:     aws.ToString(resp.Credentials.AccessKeyId),
	// 	SecretAccessKey: aws.ToString(resp.Credentials.SecretAccessKey),
	// 	SessionToken:    aws.ToString(resp.Credentials.SessionToken),
	// 	Source:          "Github",
	// 	CanExpire:       true,
	// 	Expires:         *resp.Credentials.Expiration,
	// }

	credsProvider := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		),
	)

	return credsProvider
}

func generateKubeConfig(region, eksClusterName, outputPath string) error {

	stsSvc := sts.NewFromConfig(*awsConfig)

	presignClient := sts.NewPresignClient(stsSvc, sts.WithPresignClientFromClientOptions(func(o *sts.Options) {
		o.Credentials = awsConfig.Credentials
	}))

	getCallerIdentity, err := presignClient.PresignGetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{}, func(presignOptions *sts.PresignOptions) {
		presignOptions.ClientOptions = append(presignOptions.ClientOptions, func(stsOptions *sts.Options) {
			// Add clusterId Header
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue(clusterIDHeader, eksClusterName))
			// Add back useless X-Amz-Expires query param
			stsOptions.APIOptions = append(stsOptions.APIOptions, smithyhttp.SetHeaderValue("X-Amz-Expires", requestPresignParam))
		})
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	u2, _ := url.Parse(getCallerIdentity.URL)

	req := &http.Request{
		Method: getCallerIdentity.Method,
		URL:    u2,
		Header: getCallerIdentity.SignedHeader,
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", body)

	eksSvc := eks.NewFromConfig(*awsConfig)

	res, err := eksSvc.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	})
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	certificateAuthorityData, _ := base64.StdEncoding.DecodeString(*res.Cluster.CertificateAuthority.Data)

	config := &api.Config{
		Clusters: map[string]*api.Cluster{
			*res.Cluster.Name: {
				Server:                   *res.Cluster.Endpoint,
				CertificateAuthorityData: []byte(certificateAuthorityData),
			},
		},
		Contexts: map[string]*api.Context{
			*res.Cluster.Name: {
				Cluster:   *res.Cluster.Name,
				Namespace: "default",
				AuthInfo:  *res.Cluster.Name,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			*res.Cluster.Name: {
				Token: v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(getCallerIdentity.URL)),
			},
		},
		CurrentContext: *res.Cluster.Name,
	}

	configBytes, err := clientcmd.Write(*config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(outputPath, configBytes, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}
