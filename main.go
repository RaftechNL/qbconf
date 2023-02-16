package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

func main() {
	TokenWithRoleFromArn("arn:aws:iam::123:role/AWSAdminRole", "eu-west-1", "test", "raftech-test")
}

func AssumeRole(roleArn, region, roleSessionName string) (*stscreds.AssumeRoleProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %v", err)
	}

	// Create an STS client using the default config
	stsClient := sts.NewFromConfig(cfg)

	// Create an AssumeRoleProvider that will assume the specified role
	roleProvider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = roleSessionName
	})

	return roleProvider, nil
}

func TokenWithRoleFromArn(roleArn, region, roleSessionName, eksClusterName string) error {

	provider, err := AssumeRole(roleArn, region, roleSessionName)
	if err != nil {
		fmt.Println("Failed to assume role:", err)
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Errorf("failed to load SDK config: %v", err)
		return err
	}

	cfg.Credentials = provider

	stsSvc := sts.NewFromConfig(cfg)

	presignClient := sts.NewPresignClient(stsSvc, sts.WithPresignClientFromClientOptions(func(o *sts.Options) {
		o.Credentials = cfg.Credentials
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

	eksSvc := eks.NewFromConfig(cfg)

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

	err = ioutil.WriteFile("kubeconfig.yaml", configBytes, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}
