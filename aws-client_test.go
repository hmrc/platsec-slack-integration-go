// +build aws

package platsec_slack_integration_go

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()

	if cfg.Region == "" {
		t.Errorf("failed aws configuration load")
	}
}

type mockGetParameterAPI func(ctx context.Context, params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)

func (m mockGetParameterAPI) GetParameter(ctx context.Context, params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetParameterValueFromSSM(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SSMGetParameterAPI
		name   string
		expect string
	}{
		{
			client: func(t *testing.T) SSMGetParameterAPI {
				// Convert passed in function as type mockGetParameterAPI
				return mockGetParameterAPI(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
					t.Helper()
					paramValue := "23456"

					if params.Name == nil {
						t.Error("parameter name cannot be blank")
					}

					return &ssm.GetParameterOutput{
						Parameter: &types.Parameter{Value: &paramValue},
					}, nil
				})
			},
			name:   "testParamName",
			expect: "23456",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := GetParameterValueFromSSM(ctx, tt.client(t), tt.name)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; e != *a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}