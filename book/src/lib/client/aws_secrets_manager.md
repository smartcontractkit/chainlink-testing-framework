# AWS Secrets Manager

This simple client makes it even easier to:
* read
* create
* remove secrets from AWS Secrets Manager.

Creating a new instance is straight-forward. You should either use environment variables or shared configuration and credentials.

> [!NOTE]
> Environment variables take precedence over shared credentials.

## Using environment variables
You can pass required configuration as following environment variables:
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_REGION`

## Using shared credentials
If you have shared credentials stored in `.aws/credentials` file, then the easiest way to configure the client is by setting
`AWS_PROFILE` environment variable with the profile name. If that environment variable is not set, the SDK will try to use default profile.

> [!WARNING]
> Remember, that most probably you will need to manually create a new session for that profile before running your application.


> [!NOTE]
> You can read more about configuring the AWS SDK [here](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

Once you have an instance of AWS Secrets Manager you gain access to following functions:
* `CreateSecret(key string, val string, override bool) error`
* `GetSecret(key string) (AWSSecret, error)`
* `RemoveSecret(key string, noRecovery bool) error`