package quotesdiscovery

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestQuotesAPIConsumer(t *testing.T) {

	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "QuotesDiscovery",
		Provider: "QuotesAPI",
	})
	assert.NoError(t, err)

	Convey("Given a quotes API consumer is initialized", t, func() {
		mockProvider.
			AddInteraction().
			Given("A request to fetch quotes").
			UponReceiving("A request to fetch quotes").
			WithRequest("GET", "/fetch/quotes", func(r *consumer.V4RequestBuilder) {
				r.Query("page", matchers.Integer(1))
			}).
			WillRespondWith(200, func(r *consumer.V4ResponseBuilder) {
				r.JSONBody([]Quotes{
					{
						Quote:  "Test Quote",
						Author: "Test Author",
					},
				})
			})
		Convey("When the fetch quotes function is executed", func() {
			var quotes []Quotes
			err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewQuotesAPIClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					AddSource: true,
				})))
				quotes, err = client.FetchQuotes(context.Background(), 1)
				return err
			})
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
				So(quotes, ShouldNotBeNil)
				So(len(quotes), ShouldBeGreaterThan, 0)
			})
		})
	})
}
