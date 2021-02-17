package slack

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/tomtwinkle/attendance-client/config"
)

type ClockType int

const (
	ClockTypeClockIn ClockType = iota + 1
	ClockTypeClockOut
	ClockTypeGoOut
	ClockTypeReturned
)

type SlackClient interface {
	Action(ctx context.Context, clockType ClockType) ([]*SlackResult, error)
}

type slackClient struct {
	config *config.ConfigSlack
	client *slack.Client
}

type SlackResult struct {
	ChannelId   string
	ChannelName string
	Timestamp   string
}

func NewSlackClient(cfg *config.ConfigSlack) SlackClient {
	if cfg == nil || cfg.Token == "" {
		return &slackClient{}
	}
	return &slackClient{config: cfg, client: slack.New(cfg.Token)}
}

func (s slackClient) Action(ctx context.Context, clockType ClockType) ([]*SlackResult, error) {
	if s.client == nil {
		return nil, nil
	}
	results := make([]*SlackResult, len(s.config.Channels))
	for i, channel := range s.config.Channels {
		if channel.Name == "" {
			continue
		}
		result := &SlackResult{
			ChannelName: channel.Name,
		}
		var err error
		switch clockType {
		case ClockTypeClockIn:
			if channel.ClockIn == nil {
				continue
			}
			result.ChannelId, result.Timestamp, err = s.postMessage(ctx, channel.Name, channel.ClockIn.Message)
		case ClockTypeClockOut:
			if channel.ClockOut == nil {
				continue
			}
			result.ChannelId, result.Timestamp, err = s.postMessage(ctx, channel.Name, channel.ClockOut.Message)
		case ClockTypeGoOut:
			if channel.GoOut == nil {
				continue
			}
			result.ChannelId, result.Timestamp, err = s.postMessage(ctx, channel.Name, channel.GoOut.Message)
		case ClockTypeReturned:
			if channel.Returned == nil {
				continue
			}
			result.ChannelId, result.Timestamp, err = s.postMessage(ctx, channel.Name, channel.Returned.Message)
		}
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (s slackClient) postMessage(ctx context.Context, channelName, message string) (string, string, error) {
	if channelName == "" || message == "" {
		return "", "", nil
	}
	channel, err := s.searchChannel(ctx, channelName)
	if err != nil {
		return "", "", err
	}

	// returned channelID, timestamp, err
	return s.client.PostMessageContext(
		ctx,
		channel.ID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	)
}

func (s slackClient) searchChannel(ctx context.Context, channelName string) (*slack.Channel, error) {
	var cursor string
	for {
		requestParam := &slack.GetConversationsParameters{
			Types:           []string{"public_channel", "private_channel"},
			Limit:           1000,
			ExcludeArchived: "true",
		}
		if cursor != "" {
			requestParam.Cursor = cursor
		}
		var channels []slack.Channel
		var err error
		channels, cursor, err = s.client.GetConversationsContext(ctx, requestParam)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		channelRef := &channels
		for i := 0; i < 0; i++ {
			channel := (*channelRef)[i]
			if channel.Name == channelName {
				return &channel, nil
			}
		}
		if cursor == "" {
			break
		}
	}
	return nil, errors.New(fmt.Sprintf("channel not found. channelName=%s", channelName))
}
