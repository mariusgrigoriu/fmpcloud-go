package fmpcloud

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

// HTTPClient ...
type HTTPClient struct {
	logger              *slog.Logger
	client              *resty.Client
	apiKey              string
	mainRateLimiter     *rate.Limiter
	endpointRateLimiter map[string]*rate.Limiter
	retryCount          *int
	retryWaitTime       *time.Duration
}

// Get ...
func (h *HTTPClient) Get(endpoint string, queryParams map[string]string) (response *resty.Response, err error) {
	return h.get(endpoint, queryParams, false)
}

func (h *HTTPClient) get(endpoint string, queryParams map[string]string, doNotParse bool) (response *resty.Response, err error) {
	if queryParams == nil {
		queryParams = make(map[string]string)
	}

	queryParams["apikey"] = h.apiKey

	retries := 0
	for retries < *h.retryCount {
		if retries > 0 {
			time.Sleep(*h.retryWaitTime)
		}

		endpointLimiter := h.endpointRateLimiter[endpoint]
		if endpointLimiter != nil {
			err = endpointLimiter.Wait(context.Background())
		}
		if err != nil {
			return
		}
		err = h.mainRateLimiter.Wait(context.Background())
		if err != nil {
			return
		}

		response, err = h.client.R().
			SetDoNotParseResponse(doNotParse).
			SetQueryParams(queryParams).
			Get(endpoint)

		if err != nil || response.StatusCode() != http.StatusOK {
			retries++

			// response is not valid when there is an error
			if err == nil {
				err = fmt.Errorf("http %s", response.Status())
			}

			h.logger.Info(
				"Retry request.",
				"tries", retries,
				"err", err,
				"endpoint", endpoint,
				"queryParams", queryParams,
			)

			continue
		}

		if err == nil {
			break
		}
	}

	return response, err
}

func (h *HTTPClient) EODBatchPrices(date time.Time) (response *resty.Response, err error) {
	return h.get(urlAPIStockEODBatchPrices, map[string]string{"date": date.Format("2006-01-02")}, true)
}

func (h *HTTPClient) BulkBalanceSheetStatement(year int, period string) (response *resty.Response, err error) {
	return h.get(urlAPICompanyValuationBulkBalanceSheetStatement, map[string]string{
		"year":   fmt.Sprint(year),
		"period": period,
	}, true)
}

func (h *HTTPClient) BulkCashFlowStatement(year int, period string) (response *resty.Response, err error) {
	return h.get(urlAPICompanyValuationBulkCashFlowStatement, map[string]string{
		"year":   fmt.Sprint(year),
		"period": period,
	}, true)
}

func (h *HTTPClient) BulkIncomeStatement(year int, period string) (response *resty.Response, err error) {
	return h.get(urlAPICompanyValuationBulkIncomeStatement, map[string]string{
		"year":   fmt.Sprint(year),
		"period": period,
	}, true)
}

type StatementFn func(int, string) (*resty.Response, error)
