package gcloud

import (
	"golang.org/x/net/context"

	gs "cloud.google.com/go/storage"
)

type gcloudProvider struct {}

func GCloudProvider () gcloudProvider {
	return gcloudProvider{}
}

func getContext () (context.Context) {
	return context.Background()
}

func getClient (ctx context.Context) (*gs.Client, error) {
	client, err := gs.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
