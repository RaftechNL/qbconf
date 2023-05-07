# qbconf

![Logo](https://img.raftech.nl/logo-qbconf.png)

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
* generate `<cloud>` - generates a kubeconfig file for a cluster in selected cloud provider

### generate
Generate is our root working command. It supports cloud providers ( AWS at the moment ).

#### AWS 
```
## generates kubeconfig for aws eks cluster by assuming given role ( uses provided credentials )
qbconf generate aws --cluster-name XXX --region us-east-1 --with-assume-role --role-arn "arn:aws:iam::12334556:role/AWSMagicRole"

## generates kubeconfig for aws eks cluster by assuming given role ( uses oidc credentials )
qbconf generate aws --cluster-name XXX --region us-east-1 --with-gha-oidc --role-arn "arn:aws:iam::12334556:role/AWSMagicRole"

## generates kubeconfig for aws eks cluster by using provided credentials
qbconf generate aws --cluster-name XXX --region us-east-1 
```

##### Output
The CLI will by default output a kubeconfig file called `kubeconfig.yaml`. This can be changed by using the `--output-file` flag.

## Contributing

Contributions are always welcome!


## Authors

- [@rafpe](https://www.github.com/rafpe)
