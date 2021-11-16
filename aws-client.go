package platsec_slack_integration_go

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// GetParameterValueFromSSM returns a parameter value from AWS secrets store for a given name.
func GetParameterValueFromSSM(ctx context.Context, api SSMGetParameterAPI, parameterName string) (*string, error) {
	parameterOutput, err := api.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: true,
	})
	if err != nil {
		return nil, err
	}

	return parameterOutput.Parameter.Value, nil
}

// LoadConfig creates a config to be use with aws clients.
func LoadConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	return cfg
}

// GenerateSSMClient creates the client to interact with SSM.
func GenerateSSMClient(config aws.Config) SSMGetParameterAPI {
	client := ssm.NewFromConfig(config)
	return client
}
