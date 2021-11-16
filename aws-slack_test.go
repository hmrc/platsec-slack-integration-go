// +build slack

package platsec_slack_integration_go

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	os.Setenv("SLACK_API_URL", "https://slack-notifications.tax.service.gov.uk/slack-notifications/notification")
	os.Setenv("SLACK_TOKEN_KEY", "/service_accounts/platsec_alerts_slack_password")
	os.Setenv("SLACK_USERNAME_KEY", "/service_accounts/platsec_alerts_slack_username")
	os.Setenv("SSM_READ_ROLE", "platsec_compliance_alerting_read_ssm_parameters_role")
	os.Setenv("AWS_ACCOUNT","123456789")
	extVal := m.Run()
	os.Exit(extVal)
}

func Test_createNewSlackMessage_returns_message(t *testing.T) {
	channels := []string{"TestChannel"}
	_, err := NewSlackMessage(channels, "testHeader",
		"testTile", "testMessage", "Red")
	if err != nil {
		t.Error("NewSlackMessage not created")
	}
}

func Test_createNewSlackMessage_returns_error_empty_channel(t *testing.T) {
	var channels []string
	want := "no channels specified"
	_, err := NewSlackMessage(channels, "testHeader",
		"testTile", "testMessage", "Red")

	got := err.Error()

	if got != want {
		t.Error("incorrect error returned for missing channel")
	}
}

func Test_buildHeaders_returns_valid_map(t *testing.T) {
	want := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic mteasdal:12345685",
	}

	got := buildHeaders(NewSlackNotifierConfig("mteasdal", "12344", "testUrl","",""))

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Incorrect build_headers result expected: %v got: %v", want, got)
	}
}

func Test_getEnvVariables_returns_values(t *testing.T) {
	want := map[string]string{
		"SLACK_API_URL":      "https://slack-notifications.tax.service.gov.uk/slack-notifications/notification",
		"SLACK_USERNAME_KEY": "/service_accounts/platsec_alerts_slack_username",
		"SLACK_TOKEN_KEY":    "/service_accounts/platsec_alerts_slack_password",
		"SSM_READ_ROLE":      "platsec_compliance_alerting_read_ssm_parameters_role",
		"AWS_ACCOUNT": "123456789",
	}

	got := getEnvConfig()
	t.Cleanup(func() {
		os.Unsetenv("SLACK_API_URL")
		os.Unsetenv("SLACK_USERNAME_KEY")
		os.Unsetenv("SLACK_TOKEN_KEY")
		os.Unsetenv("SSM_READ_ROLE")
		os.Unsetenv("AWS_ACCOUNT")
	})

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("config environmental variables not valid %s", diff)
	}
}

func Test_validateEnvironmentVariables_single_key_returns_false(t *testing.T) {
	testKeys := []string{"testKey1"}
	want := false
	got := validateEnvConfig(testKeys, 0)

	if want != got {
		t.Errorf("%v returned expected %v", got, want)
	}
}

func Test_createConfig_returns_valid_config_struct(t *testing.T) {

	cases :=[]struct{
		configItems map[string]string
		expected SlackNotifierConfig
	}{
		{
			configItems : map[string]string{
			"SLACK_API_URL":"testURL",
			"SLACK_USERNAME_KEY":"testUsername",
			"SLACK_TOKEN_KEY":"testToken",
			"SSM_READ_ROLE":"testRole",
			"AWS_ACCOUNT":"123456879",
			},
			expected: SlackNotifierConfig{
				username: "testUsername",
				apiUrl: "testURL",
				token: "testToken",
				ssmRole: "testRole",
				awsAccount: "123456789",
			},
		},
	}
	for _, c := range cases {
		actual:= assignConfigItems(c.configItems)
		if actual == (SlackNotifierConfig{}) {
			t.Error("empty config item created")
		}
	}
}

func Test_createConfig_returns_empty_config_struct(t *testing.T) {
	cases := []struct {
		configItems map[string]string
		expected SlackNotifierConfig
	}{
		{
			configItems : map[string]string{
				"SLACK_API_URL":"testURL",
				"SLACK_USERNAME_KEY":"testUsername",
				"SLACK_TOKEN_KEY":"testToken",
				"SSM_READ_ROLE":"testRole",
				"AWS_ACCOUNT":"123456879",
				"DUMMMY_VAR":"falseValue",
			},
			expected : SlackNotifierConfig{},
		},
		{
			configItems : map[string]string{
				"SLACK_API_URL":"testURL",
				"SSM_READ_ROLE":"testRole",
				"AWS_ACCOUNT":"123456879",
			},
			expected : SlackNotifierConfig{},
		},
	}

	for _ ,c:= range cases{
		actual := assignConfigItems(c.configItems)
		if actual != (SlackNotifierConfig{}){
			t.Error("empty config item not returned")
		}
	}
}