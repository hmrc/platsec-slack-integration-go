package platsec_slack_integration_go

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
}

// SlackService represents a struct holding the values returned by the SSM.
type SlackService struct {
	ssmUser  string
	ssmToken string
	apiUrl   string
}

func (s SlackService) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	resp, err = http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetParameterValueFromSSM returns a parameter value from AWS secrets store for a given name.
func getParameterValueFromSSM(ctx context.Context, api SSMGetParameterAPI, parameterName string) (*string, error) {
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
func loadConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	return cfg
}

// generateSSMClient creates the client to interact with SSM.
func generateSSMClient(config aws.Config) SSMGetParameterAPI {
	client := ssm.NewFromConfig(config)

	return client
}

// generateSlackService retrieves values from SSM and returns them as Service struct.
func generateSlackService(ctx context.Context, config SlackNotifierConfig, client SSMGetParameterAPI) (SlackService, error) {
	slackService := SlackService{}
	params := &ssm.GetParametersInput{
		Names:          []string{config.username, config.token},
		WithDecryption: true,
	}
	paramOutput, err := client.GetParameters(ctx, params)
	if err != nil {
		return SlackService{}, err
	}

	for _, param := range paramOutput.Parameters {
		if *param.Name == config.username {
			slackService.ssmUser = *param.Value
		} else {
			slackService.ssmToken = *param.Value
		}
	}

	return slackService, nil
}
