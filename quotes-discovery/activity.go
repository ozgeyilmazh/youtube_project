package quotesdiscovery

//go:generate mockgen -destination=activity_mock.go -package=quotesdiscovery . ActivityQuotesAPIClient,Repository
import (
	"context"
	"fmt"
)

type ActivityQuotesAPIClient interface {
	FetchQuotes(ctx context.Context, page int) ([]Quotes, error)
}

type Activity struct {
	client     ActivityQuotesAPIClient
	repository Repository
}

func NewActivity(client ActivityQuotesAPIClient, repository Repository) *Activity {
	return &Activity{client: client, repository: repository}
}

func (a *Activity) FetchQuotes(ctx context.Context, page int) ([]Quotes, error) {
	quotes, err := a.client.FetchQuotes(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quotes: %v", err)
	}

	return quotes, nil
}

func (a *Activity) BulkInsertData(ctx context.Context, quotes []Quotes) error {
	err := a.repository.BulkInsertData(ctx, quotes)
	if err != nil {
		return fmt.Errorf("failed to bulk insert data: %v", err)
	}
	return nil
}
