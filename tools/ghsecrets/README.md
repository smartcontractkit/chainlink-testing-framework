# ghsecrets

`ghsecrets` is a command-line tool designed to manage and set test secrets in either:

- **GitHub** (via the GitHub CLI), or
- **AWS Secrets Manager**.

This tool helps streamline the process of storing test secrets which can be referenced by your workflows or other services.

---

## Installation

To install the `ghsecrets` CLI, ensure you have Go installed. Then run:

```sh
go install github.com/smartcontractkit/chainlink-testing-framework/tools/ghsecrets@latest
```

Note: If you plan to set secrets in GitHub, please also install the GitHub CLI (gh).

## Usage

### 1. Setting Secrets

By default, `ghsecrets set` assumes you want to store secrets in AWS Secrets Manager, using a file from `~/.testsecrets` (if not specified). You can change the backend to GitHub, specify a custom file path, or share the AWS secret with other IAM principals. Below are common examples:

#### a) Set secrets in AWS (default)

> **⚠️ Note:** Ensure you authenticate with AWS before using the tool:
>
> ```sh
> aws sso login --profile <your-aws-sdlc-profile-with-poweruser-role>
> ```
> Use the **SDLC** profile in AWS with **PowerUserAccess** role 

This will read from `~/.testsecrets` (by default) and create/update a secret in AWS Secrets Manager:

```sh
ghsecrets set --profile <your-aws-sdlc-profile>
```

If you’d like to specify a different file:

```sh
ghsecrets set --file /path/to/mysecrets.env --profile <your-aws-sdlc-profile>
```

If you’d like to specify a custom secret name:

```sh
ghsecrets set --secret-id my-custom-secret --profile <your-aws-sdlc-profile>
```

Note: For AWS backend, the tool automatically adds the `testsecrets/` prefix if it is missing. This ensures consistency and allows GitHub Actions to access all secrets with this designated prefix.

If you’d like to share this secret with additional AWS IAM principals (e.g., a collaborator’s account):

```sh
ghsecrets set --shared-with arn:aws:iam::123456789012:role/SomeRole --profile <your-aws-sdlc-profile>
```

You can specify multiple ARNs using commas:

```sh
ghsecrets set --shared-with arn:aws:iam::123456789012:role/SomeRole,arn:aws:iam::345678901234:root --profile <your-aws-sdlc-profile>
```

#### b) Set secrets in GitHub

```sh
ghsecrets set --backend github
```

This will:
1. Read from the default file (`~/.testsecrets`) unless `--file` is specified.
2. Base64-encode the content.
3. Create/update a GitHub secret using the GitHub CLI.

### 2. Retrieving Secrets (AWS Only)

If you want to retrieve an existing secret from AWS Secrets Manager, use:

```sh
ghsecrets get --secret-id testsecrets/MySecretName --profile <your-aws-sdlc-profile>
```

By default, it tries to decode a Base64-encoded test secret. To disable decoding use `--decode false` flag:

```sh
ghsecrets get --secret-id testsecrets/MySecretName --decode false --profile <your-aws-sdlc-profile>
```

## FAQ

<details>
<summary><strong>Q: I get "command not found: ghsecrets" after installation. How do I fix this?</strong></summary>

This error typically means the directory where Go installs its binaries is not in your system’s PATH. The binaries are usually installed in `$GOPATH/bin` or `$GOBIN`.

Steps to fix:
1. If you use `asdf`, run:

    ```sh
    asdf reshim golang
    ```

2. Otherwise, add your Go bin directory to PATH manually:
    - Find your Go bin directory:

    ```sh
    echo $(go env GOPATH)/bin
    ```

    - Add it to your shell config (e.g., `~/.bashrc`, `~/.zshrc`):

    ```sh
    export PATH="$PATH:<path-to-go-bin>"
    ```

    - Reload your shell:

    ```sh
    source ~/.bashrc  # or .zshrc, etc.
    ```

3. Alternatively, run the tool using its full path without modifying PATH:

    ```sh
    $(go env GOPATH)/bin/ghsecrets set
    ```

</details>

<details>
<summary><strong>Q: What if my AWS SSO session expires?</strong></summary>

If you see errors like `InvalidGrantException` when setting or retrieving secrets from AWS, your SSO session may have expired. Re-authenticate using:

```sh
aws sso login --profile <my-aws-profile>
```

Then try running `ghsecrets` again.

</details>

<details>
<summary><strong>Q: What if I get an error that says "GitHub CLI not found"?</strong></summary>

For GitHub secrets, this tool requires the GitHub CLI. Please install it first:

```sh
brew install gh
# or
sudo apt-get install gh
```

Then run:

```sh
gh auth login
```

and follow the prompts to authenticate.

</details>

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. Feel free to use, modify, and distribute it as needed.

