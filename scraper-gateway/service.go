package scraper_gateway

//go:generate mockgen -destination=service_mock.go -package=scraper_gateway . ServiceQuotesDiscoveryClient
import "context"

type ServiceQuotesDiscoveryClient interface {
	TriggerScrape(ctx context.Context, start int, end int) error
}

type Service struct {
	triggerScrapeClient ServiceQuotesDiscoveryClient
}

func NewService(triggerScrapeClient ServiceQuotesDiscoveryClient) *Service {
	return &Service{triggerScrapeClient: triggerScrapeClient}
}

func (s *Service) TriggerScrape(ctx context.Context, start int, end int) error {
	return nil
}
