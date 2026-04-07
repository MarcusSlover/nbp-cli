// Package nbp handles communication with the National Bank of Poland (NBP) API.
package nbp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Rate represents a single exchange rate entry from the NBP API.
type Rate struct {
	No            string  `json:"no"`            // Table number
	EffectiveDate string  `json:"effectiveDate"` // Date the rate was published
	Mid           float64 `json:"mid"`           // Middle exchange rate value
}

// CurrencyResponse represents the full response from the NBP API for a single currency.
type CurrencyResponse struct {
	Code  string `json:"code"`  // ISO 4217 currency code
	Rates []Rate `json:"rates"` // List of rates returned by the API
}

// TableEntry represents a single currency in the NBP exchange rate tables.
type TableEntry struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

// TableResponse represents the response from the NBP API for a full table.
type TableResponse struct {
	Table         string       `json:"table"`
	No            string       `json:"no"`
	EffectiveDate string       `json:"effectiveDate"`
	Rates         []TableEntry `json:"rates"`
}

// FetchNormalRate fetches the rate for a specific date or the latest one if "today" is used.
func FetchNormalRate(currency, inputDate string) (*CurrencyResponse, error) {
	url := fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/a/%s/%s/?format=json", currency, inputDate)

	if inputDate == "today" {
		url = fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/a/%s/last/1/?format=json", currency)
	}

	return callNBP[CurrencyResponse](url)
}

// FetchTaxRate fetches the rate for the last business day before targetDate.
func FetchTaxRate(currency string, targetDate time.Time) (*CurrencyResponse, error) {
	// Subtract 1 day per Polish tax law
	endDate := targetDate.AddDate(0, 0, -1)
	// Look back 7 days to ensure we hit a business day
	startDate := endDate.AddDate(0, 0, -7)

	url := fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/a/%s/%s/%s/?format=json",
		currency, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	return callNBP[CurrencyResponse](url)
}

// FetchAvailableCurrencies fetches all current currency codes from NBP (Table A).
func FetchAvailableCurrencies() ([]string, error) {
	url := "https://api.nbp.pl/api/exchangerates/tables/a/?format=json"
	tables, err := callNBP[[]TableResponse](url)
	if err != nil {
		return nil, err
	}

	if len(*tables) == 0 {
		return nil, fmt.Errorf("no tables returned from NBP")
	}

	var codes []string
	for _, rate := range (*tables)[0].Rates {
		codes = append(codes, rate.Code)
	}
	return codes, nil
}

// callNBP executes a GET request to the specified NBP API URL.
// It includes a custom User-Agent header to avoid being blocked.
func callNBP[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "NBP-CLI/1.0 (Go-Agent)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NBP API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data T
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}
