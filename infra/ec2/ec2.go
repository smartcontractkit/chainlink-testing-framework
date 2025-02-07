package main

import (
	"ec2/common"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
	"github.com/rs/zerolog/log"
)

func main() {
	defer jsii.Close()
	err := common.NewVM("self-serve", &common.VMConfig{
		Name:   "bcm-1",
		Region: "us-east-1",
		// us-east-1 latest Ubuntu
		AMI: "ami-04b4f1a9cf54c11d0",
		// eu-north-1 latest Ubuntu
		//AMI:    "ami-08eb150f611ca277f",
		Class:    awsec2.InstanceClass_T3,
		Size:     awsec2.InstanceSize_MEDIUM,
		Tags:     "cost-center=bcm,environment=test",
		UserData: common.DefaultDockerScript,
		ChiselPorts: []string{
			"8080",
		},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
