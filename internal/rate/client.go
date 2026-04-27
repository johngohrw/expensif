package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://api.frankfurter.dev/v1",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type latestResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// Latest fetches the latest exchange rates from the given base currency.
// Returns a map of currency code -> rate, the rate date, and any error.
func (c *Client) Latest(ctx context.Context, base string) (map[string]float64, string, error) {
	url := fmt.Sprintf("%s/latest?from=%s", c.baseURL, base)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("frankfurter returned %d", resp.StatusCode)
	}

	var result latestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("decode rates: %w", err)
	}

	// Include the base currency at 1.0
	result.Rates[result.Base] = 1.0

	return result.Rates, result.Date, nil
}
