package fmpcloud

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// HTTPClient ...
type HTTPClient struct {
	logger              *zap.Logger
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
				zap.Int("retries", retries),
				zap.Error(err),
				zap.String("endpoint", endpoint),
				zap.Any("data", queryParams),
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
