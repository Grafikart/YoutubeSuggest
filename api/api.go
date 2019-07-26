package api

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type API struct {
	service *youtube.Service
}

type Subscriptions []Subscription

type Subscription struct {
	ID   string
	Name string
}

func NewAPI(token string) (*API, error) {
	ctx := context.Background()
	s, err := youtube.NewService(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, fmt.Errorf("cannot create youtube service %v", err)
	}
	return &API{s}, nil
}

func (api *API) Subscriptions(channelID string) (subscriptions Subscriptions, err error) {
	ctx := context.Background()
	err = api.service.Subscriptions.
		List("snippet").
		ChannelId(channelID).
		MaxResults(50).
		Pages(ctx, appendSubscription(&subscriptions))
	if err != nil {
		return nil, fmt.Errorf("can't retrieve subscriptions for %q, %v", channelID, err)
	}
	return subscriptions, nil
}

func appendSubscription(subs *Subscriptions) func(res *youtube.SubscriptionListResponse) error {
	return func(res *youtube.SubscriptionListResponse) error {
		newSubs := *subs
		for _, s := range res.Items {
			newSubs = append(newSubs, Subscription{s.Snippet.ResourceId.ChannelId, s.Snippet.Title})
		}
		*subs = newSubs
		return nil
	}
}
