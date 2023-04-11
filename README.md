# qbconf

![Logo](https://img.raftech.nl/white_logo_color1_background.png)

A minimalistic CLI to generate kubeconfig 

#
[![License](https://img.shields.io/github/license/raftechnl/terrafile)](./LICENSE)


## Functionality

Minimalistic Kubernetes kubeconfig file generator using AWS STS and EKS APIs. It supports role assumption and Github Actions OIDC out of the box! 

Its small footprint of 4MBs and single responsibility makes it ideal for use in CI/CD pipelines.

## Installing

### Download
> Check our release page to download a specific version

```shell
    #!/bin/bash

    # Fetch the latest release version from Github API
    VERSION=$(curl --silent "https://api.github.com/repos/RaftechNL/qbconf/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    # Set the URL of the tarball for the latest release
    URL="https://github.com/RaftechNL/cli-yaml-merger/releases/download/${VERSION}/qbconf_${VERSION}_darwin_x86_64.tar.gz"

    # Download and install the latest release
    curl -L ${URL} | tar xz
    chmod +x qbconf
    sudo mv qbconf /usr/local/bin/
```

### Homebrew
```shell
brew tap  RaftechNL/toolbox
brew install raftechnl/toolbox/qbconf
```

## Usage
CLI supports the following actions
* generate ( using assume role operations )
* generate-gha ( using github actions OIDC )
  

### generate

Generating a kubeconfig file for a cluster in AWS using assume role operations requires the following information:
* eks-cluster-name
* role-arn

```shell
USAGE:
   qbconf generate [command options] [arguments...]

OPTIONS:
   --role-arn value           ARN of the AWS IAM role to assume [$AWS_ROLE_ARN]
   --region value             AWS region (default: "eu-west-1") [$AWS_REGION]
   --role-session-name value  Name of the AWS STS role session to create (default: "qbconf-session") [$AWS_ROLE_SESSION_NAME]
   --eks-cluster-name value   Name of the EKS cluster to generate a kubeconfig file for
   --output-file value        Name of the file to write the generated kubeconfig to (default: "kubeconfig.yaml")
   --help, -h                 show help
```

### generate-gha

> This action is only available when qbconf is run from runner in hosted Github Actions

Generating a kubeconfig file for a cluster in AWS using assume role with OIDC in Github Actions requires the following information:
* eks-cluster-name
* role-arn

```shell
USAGE:
   qbconf generate-gha [command options] [arguments...]

OPTIONS:
   --role-arn value           ARN of the AWS IAM role to assume [$AWS_ROLE_ARN]
   --region value             AWS region (default: "eu-west-1") [$AWS_REGION]
   --role-session-name value  Name of the AWS STS role session to create (default: "qbconf-session") [$AWS_ROLE_SESSION_NAME]
   --eks-cluster-name value   Name of the EKS cluster to generate a kubeconfig file for
   --output-file value        Name of the file to write the generated kubeconfig to (default: "kubeconfig.yaml")
   --help, -h                 show help
```

## Contributing

Contributions are always welcome!


## Authors

- [@rafpe](https://www.github.com/rafpe)
