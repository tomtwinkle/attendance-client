package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_readConfig(t *testing.T) {
	testConfig := &config{ConfigPath: "testdata/fixture.yaml"}
	actual, err := testConfig.readConfig()
	assert.NoError(t, err)
	assert.Equal(t, &ConfigYaml{
		TenantCode: "test_tenant",
		OBCiD:      "100",
		Password:   "password",
		Slack: &ConfigSlack{
			Token: "test_slack_token",
			Channels: []ConfigSlackChannel{
				{
					Name:     "channelA",
					ClockIn:  &ConfigSlackAction{Message: "channelA post clockin"},
					ClockOut: &ConfigSlackAction{Message: "channelA post clockout"},
					GoOut:    &ConfigSlackAction{Message: "channelA post goout"},
					Returned: &ConfigSlackAction{Message: "channelA post returned"},
				},
				{
					Name:     "channelB",
					ClockIn:  &ConfigSlackAction{Message: "channelB post clockin"},
					ClockOut: &ConfigSlackAction{Message: "channelB post clockout"},
					GoOut:    &ConfigSlackAction{Message: "channelB post goout"},
					Returned: &ConfigSlackAction{Message: "channelB post returned"},
				},
			},
		},
	}, actual)
}
