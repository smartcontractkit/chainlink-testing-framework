module github.com/smartcontractkit/chainlink-testing-framework/tools/ghsecrets

go 1.24.0

require (
	github.com/aws/aws-sdk-go-v2 v1.31.0
	github.com/aws/aws-sdk-go-v2/config v1.27.39
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.33.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.31.3
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.37 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.23.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.27.3 // indirect
	github.com/aws/smithy-go v1.21.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
