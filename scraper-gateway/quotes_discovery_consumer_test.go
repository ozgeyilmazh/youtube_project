package scraper_gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestQuotesDiscoveryConsumer(t *testing.T) {

	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "QuotesScraperGateway",
		Provider: "QuotesDiscovery",
	})
	assert.NoError(t, err)
	Convey("Given a quotes discovery consumer is initialized", t, func() {

		mockProvider.
			AddInteraction().
			Given("A request to trigger scraping").
			UponReceiving("A request to trigger scraping").
			WithRequest("POST", "/scrape/quotes/start/1/end/10").
			WillRespondWith(202)

		Convey("When the trşşger scrape function is executed", func() {

			err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewQuotesDiscoveryClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
				err := client.TriggerScrape(context.Background(), 1, 10)
				return err
			})
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
