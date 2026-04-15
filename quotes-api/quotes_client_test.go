package quotesapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestQuotesClient(t *testing.T) {

	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "quotes-api-client",
		Provider: "quotes-api",
		Host:     "127.0.0.1",
		TLS:      false,
	})
	if err != nil {
		t.Fatalf("Failed to create mock provider: %v", err)
	}

	t.Run("Given valid page When FetchQuotes is called Then it should return the quotes", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("A request to fetch quotes").
			UponReceiving("A request to fetch quotes").
			WithRequest("GET", "/api/quotes", func(r *consumer.V4RequestBuilder) {
				r.Query("page", matchers.Integer(1))
			}).
			WillRespondWith(200, func(r *consumer.V4ResponseBuilder) {
				r.JSONBody([]QuotesResponse{
					{Quote: "Test Quote", Author: "Test Author"},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				QUOTES_API_URL = fmt.Sprintf("http://%s:%d/api/quotes?page=", config.Host, config.Port)
				client, err := NewClient(getLogger())
				if err != nil {
					return err
				}

				quotes, err := client.FetchQuotes(context.Background(), 1)
				if err != nil {
					return err
				}
				assert.NoError(t, err)
				assert.Equal(t, 1, len(quotes))
				assert.Equal(t, "Test Quote", quotes[0].Quote)
				assert.Equal(t, "Test Author", quotes[0].Author)
				return nil
			})

		if err != nil {
			t.Fatalf("Failed to fetch quotes: %v", err)
		}
	})
}
