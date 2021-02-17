package cli

import (
	"context"
	"github.com/tomtwinkle/attendance-client/config"
	"github.com/tomtwinkle/attendance-client/slack"
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"log"
)

type cli struct {
	bugyoClient bugyoclient.BugyoClient
	slackClient slack.SlackClient
}

type CLI interface {
	PunchMark(ctx context.Context, clockType bugyoclient.ClockType, slackOnly bool) error
}

func NewCLI() CLI {
	cfg := config.NewConfig()
	bCfg, err := cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
	c, err := bugyoclient.NewClient(&bugyoclient.BugyoConfig{
		TenantCode: bCfg.TenantCode,
		OBCiD:      bCfg.OBCiD,
		Password:   bCfg.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	if bCfg.Slack == nil {
		return &cli{c, nil}
	}
	return &cli{c, slack.NewSlackClient(bCfg.Slack)}
}

func (c cli) PunchMark(ctx context.Context, clockType bugyoclient.ClockType, slackOnly bool) error {
	if !slackOnly {
		if err := c.bugyoClient.Login(); err != nil {
			return err
		}
		if err := c.bugyoClient.Punchmark(clockType); err != nil {
			return err
		}
		log.Printf("success Punchmark [%s]", clockType)
	}

	if c.slackClient != nil {
		var results []*slack.SlackResult
		var err error
		switch clockType {
		case bugyoclient.ClockTypeClockIn:
			results, err = c.slackClient.Action(ctx, slack.ClockTypeClockIn)
		case bugyoclient.ClockTypeClockOut:
		case bugyoclient.ClockTypeGoOut:
		case bugyoclient.ClockTypeReturned:
		}
		if err != nil {
			return err
		}
		for _, r := range results {
			log.Printf("success Slack Action [%s] %s - %s", clockType, r.ChannelName, r.Timestamp)
		}
	} else if slackOnly {
		log.Print("use slackOnly option, but slack config is not set")
	}
	return nil
}
