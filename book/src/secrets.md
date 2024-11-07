## Using AWSSecretsManager from code

`client/secretsmanager.go` has a simple API to read/write/delete secrets.

It uses a struct to protect such secrets from accidental printing or marshalling, see an [example](../../lib/client/secretsmanager_test.go) test

## Using AWSSecretsManager via CLI

To create a static secret use `aws cli`

```
aws --region us-west-2 secretsmanager create-secret \
    --name MyTestSecret \
    --description "My test secret created with the CLI." \
    --secret-string "{\"user\":\"diegor\",\"password\":\"EXAMPLE-PASSWORD\"}"
```

Example of reading the secret

```
aws --region us-west-2 secretsmanager get-secret-value --secret-id MyTestSecret
```

For more information check [AWS CLI Reference](https://docs.aws.amazon.com/cli/v1/userguide/cli_secrets-manager_code_examples.html)
