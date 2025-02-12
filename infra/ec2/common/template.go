package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"text/template"
)

const (
	DefaultDockerScript = `#!/bin/bash
sudo apt update -y
sudo apt upgrade -y
sudo apt install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update -y
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo usermod -aG docker ubuntu
sudo systemctl enable docker
sudo systemctl start docker
curl https://i.jpillora.com/chisel! | bash
`
	CRUNDockerScript = `#!/bin/bash
sudo apt update -y
sudo apt upgrade -y
sudo apt install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update -y
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo usermod -aG docker ubuntu
sudo systemctl enable docker
sudo systemctl start docker

# CRUN Docker runtime

set -e

# Install dependencies
sudo apt update -y
sudo apt-get install -y make git gcc build-essential pkgconf libtool \
   libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev \
   go-md2man autoconf python3 automake

# Clone and build crun
git clone https://github.com/containers/crun.git
cd crun
./autogen.sh
./configure
make
sudo make install

# Configure Docker to use crun
sudo mkdir -p /etc/docker
echo '{
  "runtimes": {
    "crun": {
      "path": "/usr/local/bin/crun"
    }
  },
  "default-runtime": "crun"
}' | sudo tee /etc/docker/daemon.json

# Restart Docker to apply changes
sudo systemctl restart docker

# Verify crun is the default runtime
sudo docker info | grep "Default Runtime"
`
)

// VMConfig defines the configuration for the VM
type VMConfig struct {
	Region      string
	AMI         string
	Class       awsec2.InstanceClass
	Size        awsec2.InstanceSize
	UserData    string
	Name        string
	ChiselPorts []string
	Tags        string
}

// NewVM creates a new EC2 instance with the given configuration
func NewVM(stackName string, config *VMConfig) error {
	app := awscdk.NewApp(nil)

	action := app.Node().TryGetContext(jsii.String("action"))
	if action == nil {
		log.Fatal().Msg("action should be used, cdk $cmd -c action=...")
		os.Exit(1)
	}
	a := action.(string)
	log.Info().Str("Action", a).Msg("Running action")
	switch a {
	case "bootstrap":
		return nil
	case "deploy":
		// Create a new stack
		stack := awscdk.NewStack(app, jsii.String(stackName), nil)

		// Generate and import key pair
		keyPairName, err := generateAndImportKeyPair(config)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to generate and import key pair")
		}

		// VPC
		vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{
			MaxAzs: jsii.Number(1),
		})

		// Hardcoded Security Group
		sg := awsec2.NewSecurityGroup(stack, jsii.String("SecurityGroup"), &awsec2.SecurityGroupProps{
			Vpc:              vpc,
			AllowAllOutbound: jsii.Bool(true),
			Description:      jsii.String("Allow ports 8080 and 8082"),
		})
		sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("Allow SSH"), nil)
		// tailscale or chisel port
		sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(44044)), jsii.String("Allow HTTP"), nil)

		// Hardcoded EC2 Instance Role
		role := awsiam.NewRole(stack, jsii.String("InstanceRole"), &awsiam.RoleProps{
			AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		})

		// Parse tags
		tags := parseTags(config.Tags)

		// Create EC2 instance
		inst := awsec2.NewInstance(stack, jsii.String(config.Name), &awsec2.InstanceProps{
			Vpc:          vpc,
			InstanceType: awsec2.InstanceType_Of(config.Class, config.Size),
			MachineImage: awsec2.NewGenericLinuxImage(&map[string]*string{
				config.Region: aws.String(config.AMI),
			}, nil),
			KeyName:       jsii.String(keyPairName),
			VpcSubnets:    &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
			SecurityGroup: sg,
			Role:          role,
			UserData:      awsec2.UserData_Custom(jsii.String(config.UserData)),
		})

		// Add tags to the instance
		for key, value := range tags {
			awscdk.Tags_Of(inst).Add(jsii.String(key), jsii.String(value), nil)
		}
		app.Synth(nil)
		if err := GenerateConnectScript(stackName, config); err != nil {
			return err
		}
		if err := GenerateForwardScript(stackName, config); err != nil {
			return err
		}
	case "destroy":
	default:
		log.Fatal().Str("Action", a).Msg("unsupported action")
	}
	return nil
}

// generateAndImportKeyPair generates an RSA key pair and imports it into AWS
func generateAndImportKeyPair(config *VMConfig) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate RSA key: %w", err)
	}
	privateKeyFile := "cdk-ec2-keypair.pem"
	file, err := os.Create(privateKeyFile)
	if err != nil {
		return "", fmt.Errorf("failed to create key file: %w", err)
	}
	defer file.Close()

	_ = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	_ = os.Chmod(privateKeyFile, 0400)
	fmt.Printf("Private key saved to %s\n", privateKeyFile)
	publicKeyBytes, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH public key: %w", err)
	}
	sshPublicKey := string(ssh.MarshalAuthorizedKey(publicKeyBytes))

	// import public key
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	}))
	svc := ec2.New(sess)
	keyPairName := fmt.Sprintf("ec2-dev-key-%s", uuid.NewString()[0:4])
	_, err = svc.ImportKeyPair(&ec2.ImportKeyPairInput{
		KeyName:           aws.String(keyPairName),
		PublicKeyMaterial: []byte(sshPublicKey),
	})
	if err != nil {
		return "", fmt.Errorf("failed to import key pair: %w", err)
	}
	log.Info().Str("Key", keyPairName).Msg("Key pair imported successfully")
	return keyPairName, nil
}

// parseTags parses a comma-separated string of key=value pairs into a map
func parseTags(tags string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(tags, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

// GenerateConnectScript generates a connect.sh script that is used to connect via SSH
func GenerateConnectScript(stackName string, config *VMConfig) error {
	// Define the template for the connect.sh script
	const connectTemplate = `#!/bin/bash
set -e

# Retrieve the public IP of the instance
INSTANCE_PUBLIC_IP=$(aws ec2 describe-instances \
    --filters "Name=tag:Name,Values={{.StackName}}/{{.VMName}}" \
    --query "Reservations[*].Instances[*].PublicIpAddress" \
    --output text)

echo "Instance Public IP: $INSTANCE_PUBLIC_IP"

# SSH into the instance with port forwarding
ssh -o StrictHostKeyChecking=no -i cdk-ec2-keypair.pem ubuntu@$INSTANCE_PUBLIC_IP
`

	// Parse the template
	tmpl, err := template.New("connect.sh").Parse(connectTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the connect.sh file
	file, err := os.Create("connect.sh")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Prepare the data for the template
	data := struct {
		StackName string
		VMName    string
	}{
		StackName: stackName,
		VMName:    config.Name,
	}

	// Execute the template and write to the file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Make the script executable
	err = os.Chmod("connect.sh", 0755)
	if err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	return nil
}

// GenerateForwardScript generates a forward.sh script with optional port forwarding
func GenerateForwardScript(stackName string, config *VMConfig) error {
	// Define the template for the connect.sh script
	const connectTemplate = `#!/bin/bash
set -e

# Retrieve the public IP of the instance
INSTANCE_PUBLIC_IP=$(aws ec2 describe-instances \
    --filters "Name=tag:Name,Values={{.StackName}}/{{.VMName}}" \
    --query "Reservations[*].Instances[*].PublicIpAddress" \
    --output text)

echo "Instance Public IP: $INSTANCE_PUBLIC_IP"

# Forward multiple ports
chisel client --fingerprint $CHISEL_FINGERPRINT $INSTANCE_PUBLIC_IP:44044 {{range .ChiselPorts}}localhost:{{.}} {{end}}
`

	// Parse the template
	tmpl, err := template.New("connect.sh").Parse(connectTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the connect.sh file
	file, err := os.Create("forward.sh")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Prepare the data for the template
	data := struct {
		StackName   string
		VMName      string
		ChiselPorts []string
	}{
		StackName:   stackName,
		VMName:      config.Name,
		ChiselPorts: config.ChiselPorts,
	}

	// Execute the template and write to the file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Make the script executable
	err = os.Chmod("forward.sh", 0755)
	if err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	return nil
}
