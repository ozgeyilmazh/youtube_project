package quotesdiscovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPQuotesAPIClient loads quotes from the Quotable public API (https://github.com/lukePeavey/quotable).
type HTTPQuotesAPIClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPQuotesAPIClient(baseURL string) *HTTPQuotesAPIClient {
	if baseURL == "" {
		baseURL = "https://api.quotable.io"
	}
	return &HTTPQuotesAPIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type quotableQuotesResponse struct {
	Results []struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	} `json:"results"`
}

func (c *HTTPQuotesAPIClient) FetchQuotes(ctx context.Context, page int) ([]Quotes, error) {
	url := fmt.Sprintf("%s/quotes?page=%d&limit=20", c.baseURL, page)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quotes API: %s", resp.Status)
	}
	var body quotableQuotesResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	out := make([]Quotes, 0, len(body.Results))
	for _, r := range body.Results {
		out = append(out, Quotes{Quote: r.Content, Author: r.Author})
	}
	return out, nil
}
