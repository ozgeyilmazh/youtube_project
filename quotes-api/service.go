package quotesapi

//go:generate mockgen -destination=service_mock.go -package=quotesapi . QuotesAPIClient
import "context"

type Service struct {
	quotesAPIClient QuotesAPIClient
}

type QuotesAPIClient interface {
	FetchQuotes(ctx context.Context, page int) ([]QuotesResponse, error)
}

func NewService(quotesAPIClient QuotesAPIClient) *Service {
	return &Service{quotesAPIClient: quotesAPIClient}
}

func (s *Service) FetchQuotes(ctx context.Context, page int) ([]QuotesResponse, error) {
	return s.quotesAPIClient.FetchQuotes(ctx, page)
}
