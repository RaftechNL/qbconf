# v1.27.3 (2023-02-15)

* **Announcement**: When receiving an error response in restJson-based services, an incorrect error type may have been returned based on the content of the response. This has been fixed via PR #2012 tracked in issue #1910.
* **Bug Fix**: Correct error type parsing for restJson services.

# v1.27.2 (2023-02-08)

* No change notes available for this release.

# v1.27.1 (2023-02-03)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.27.0 (2023-01-05)

* **Feature**: Add `ErrorCodeOverride` field to all error structs (aws/smithy-go#401).

# v1.26.0 (2022-12-15)

* **Feature**: Add support for Windows managed nodes groups.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.25.0 (2022-12-07)

* **Feature**: Adds support for EKS add-ons configurationValues fields and DescribeAddonConfiguration function

# v1.24.1 (2022-12-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.24.0 (2022-11-29)

* **Feature**: Adds support for additional EKS add-ons metadata and filtering fields

# v1.23.0 (2022-11-16)

* **Feature**: Adds support for customer-provided placement groups for Kubernetes control plane instances when creating local EKS clusters on Outposts

# v1.22.4 (2022-11-07)

* No change notes available for this release.

# v1.22.3 (2022-10-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.2 (2022-10-21)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.1 (2022-09-20)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.22.0 (2022-09-14)

* **Feature**: Adding support for local Amazon EKS clusters on Outposts
* **Feature**: Adds support for EKS Addons ResolveConflicts "preserve" flag. Also adds new update failed status for EKS Addons.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.11 (2022-09-02)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.10 (2022-08-31)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.9 (2022-08-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.8 (2022-08-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.7 (2022-08-09)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.6 (2022-08-08)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.5 (2022-08-01)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.4 (2022-07-05)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.3 (2022-06-29)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.2 (2022-06-07)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.1 (2022-05-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.21.0 (2022-05-10)

* **Feature**: Adds BOTTLEROCKET_ARM_64_NVIDIA and BOTTLEROCKET_x86_64_NVIDIA AMI types to EKS managed nodegroups

# v1.20.7 (2022-05-03)

* No change notes available for this release.

# v1.20.6 (2022-04-25)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.5 (2022-04-12)

* No change notes available for this release.

# v1.20.4 (2022-03-30)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.3 (2022-03-24)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.2 (2022-03-23)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.20.1 (2022-03-08.3)

* No change notes available for this release.

# v1.20.0 (2022-03-08)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.19.0 (2022-02-24)

* **Feature**: API client updated
* **Feature**: Adds RetryMaxAttempts and RetryMod to API client Options. This allows the API clients' default Retryer to be configured from the shared configuration files or environment variables. Adding a new Retry mode of `Adaptive`. `Adaptive` retry mode is an experimental mode, adding client rate limiting when throttles reponses are received from an API. See [retry.AdaptiveMode](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/retry#AdaptiveMode) for more details, and configuration options.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.18.0 (2022-01-14)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.17.0 (2022-01-07)

* **Feature**: API client updated
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.16.0 (2021-12-21)

* **Feature**: API Paginators now support specifying the initial starting token, and support stopping on empty string tokens.

# v1.15.1 (2021-12-02)

* **Bug Fix**: Fixes a bug that prevented aws.EndpointResolverWithOptions from being used by the service client. ([#1514](https://github.com/aws/aws-sdk-go-v2/pull/1514))
* **Dependency Update**: Updated to the latest SDK module versions

# v1.15.0 (2021-11-30)

* **Feature**: API client updated

# v1.14.0 (2021-11-19)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.13.0 (2021-11-12)

* **Feature**: Waiters now have a `WaitForOutput` method, which can be used to retrieve the output of the successful wait operation. Thank you to [Andrew Haines](https://github.com/haines) for contributing this feature.

# v1.12.0 (2021-11-06)

* **Feature**: The SDK now supports configuration of FIPS and DualStack endpoints using environment variables, shared configuration, or programmatically.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Feature**: Updated service to latest API model.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.11.0 (2021-10-21)

* **Feature**: Updated  to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.2 (2021-10-11)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.1 (2021-09-17)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.10.0 (2021-09-10)

* **Feature**: API client updated

# v1.9.0 (2021-08-27)

* **Feature**: Updated API model to latest revision.
* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.2 (2021-08-19)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.1 (2021-08-04)

* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.8.0 (2021-07-15)

* **Feature**: Updated service model to latest version.
* **Dependency Update**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.7.0 (2021-06-25)

* **Feature**: Updated `github.com/aws/smithy-go` to latest version
* **Dependency Update**: Updated to the latest SDK module versions

# v1.6.0 (2021-06-11)

* **Feature**: Updated to latest API model.

# v1.5.0 (2021-05-20)

* **Feature**: API client updated
* **Dependency Update**: Updated to the latest SDK module versions

# v1.4.0 (2021-05-14)

* **Feature**: Constant has been added to modules to enable runtime version inspection for reporting.
* **Feature**: Updated to latest service API model.
* **Dependency Update**: Updated to the latest SDK module versions
