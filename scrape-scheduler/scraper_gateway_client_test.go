package scrape_scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestScraperGatewayClient(t *testing.T) {
	start := 1
	end := 10

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "QuotesScrapingScheduler",
		Provider: "QuotesScrapingGateway",
	})

	assert.NoError(t, err)

	Convey("Given QuotesScraperGatewya is initialized", t, func() {
		mockProvider.
			AddInteraction().
			Given("A request to trigger scraping").
			UponReceiving("A request to trigger scraping").
			WithRequest("POST", "/scrape/quotes/start/1/end/10").
			WillRespondWith(202)
		Convey("When TriggerScraping is called", func() {
			err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewQuotesScraperGatewayClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), logger)
				err := client.TriggerScraping(context.Background(), start, end)
				return err
			})
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
