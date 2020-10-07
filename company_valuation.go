package fmpcloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spacecodewor/fmpcloud-go/objects"
)

// Url const for request
const (
	urlAPICompanyValuationRSSFeed                          = "/rss_feed"
	urlAPICompanyValuationEarningCalendar                  = "/earning_calendar"
	urlAPICompanyValuationEarningsSurpises                 = "/earnings-surpises/%s"
	urlAPICompanyValuationEarningCallTranscript            = "/earning_call_transcript/%s"
	urlAPICompanyValuationHistoryEarningCalendar           = "/historical/earning_calendar/%s"
	urlAPICompanyValuationIPOCalendar                      = "/ipo_calendar"
	urlAPICompanyValuationSplitCalendar                    = "/stock_split_calendar"
	urlAPICompanyValuationDividendCalendar                 = "/stock_dividend_calendar"
	urlAPICompanyValuationInstituionalHolder               = "/institutional-holder/%s"
	urlAPICompanyValuationMutualFundHolder                 = "/mutual-fund-holder/%s"
	urlAPICompanyValuationETFHolder                        = "/etf-holder/%s"
	urlAPICompanyValuationETFSectorWeightings              = "/etf-sector-weightings/%s"
	urlAPICompanyValuationETFCountryWeightings             = "/etf-country-weightings/%s"
	urlAPICompanyValuationIncomeStatement                  = "/income-statement/%s"
	urlAPICompanyValuationIncomeStatementGrowth            = "/income-statement-growth/%s"
	urlAPICompanyValuationBalanceSheetStatement            = "/balance-sheet-statement/%s"
	urlAPICompanyValuationBalanceSheetStatementGrowth      = "/balance-sheet-statement-growth/%s"
	urlAPICompanyValuationCashFlowStatement                = "/cash-flow-statement/%s"
	urlAPICompanyValuationCashFlowStatementGrowth          = "/cash-flow-statement-growth/%s"
	urlAPICompanyValuationIncomeStatementAsReported        = "/income-statement-as-reported/%s"
	urlAPICompanyValuationBalanceSheetStatementAsReported  = "/balance-sheet-statement-as-reported/%s"
	urlAPICompanyValuationCashFlowStatementAsReported      = "/cash-flow-statement-as-reported/%s"
	urlAPICompanyValuationFinancialStatementFullAsReported = "/financial-statement-full-as-reported/%s"
	urlAPICompanyValuationFinancialRatios                  = "/ratios/%s"
	urlAPICompanyValuationFinancialRatiosTTM               = "/ratios-ttm/%s"
	urlAPICompanyValuationKeyMetrics                       = "/key-metrics/%s"
	urlAPICompanyValuationKeyMetricsTTM                    = "/key-metrics-ttm/%s"
	urlAPICompanyValuationEnterpriseValues                 = "/enterprise-values/%s"
	urlAPICompanyValuationFinancialGrowth                  = "/financial-growth/%s"
	urlAPICompanyValuationDiscountedCashFlow               = "/discounted-cash-flow/%s"
	urlAPICompanyValuationHistoryDailyDiscountedCashFlow   = "/historical-daily-discounted-cash-flow/%s"
	urlAPICompanyValuationHistoryDiscountedCashFlow        = "/historical-discounted-cash-flow-statement/%s"
	urlAPICompanyValuationRating                           = "/rating/%s"
	urlAPICompanyValuationHistoryRating                    = "/historical-rating/%s"
	urlAPICompanyValuationMarketCapitalization             = "/market-capitalization/%s"
	urlAPICompanyValuationHistoryMarketCapitalization      = "/historical-market-capitalization/%s"
	urlAPICompanyValuationDelistedCompanyList              = "/delisted-companies"
	urlAPICompanyValuationStockNews                        = "/stock_news"
	urlAPICompanyValuationStockScreener                    = "/stock-screener"
	urlAPICompanyValuationAnalystEstimates                 = "/analyst-estimates/%s"
	urlAPICompanyValuationAnalystStockRecommendations      = "/analyst-stock-recommendations/%s"
	urlAPICompanyValuationGrade                            = "/grade/%s"
	urlAPICompanyValuationPressReleases                    = "/press-releases/%s"
	urlAPICompanyValuationFinancialStatementsList          = "/financial-statement-symbol-lists"
	urlAPICompanyValuationEconomicCalendarEventList        = "/economic_calendar_event_list"
	urlAPICompanyValuationEconomicCalendar                 = "/economic_calendar"
	urlAPICompanyValuationHistoryEconomicCalendar          = "/historical/economic_calendar/%s"
)

// CompanyValuation client
type CompanyValuation struct {
	Client *resty.Client
	url    string
	apiKey string
}

// RssFeed - SEC RSS feeds is a very helpful resource for staying current on the most recent financial statements posted on the SEC
func (c *CompanyValuation) RssFeed() (fList []objects.RssFeed, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + urlAPICompanyValuationRSSFeed)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &fList)
	if err != nil {
		return nil, err
	}

	return fList, nil
}

// EarningCalendar - earning Calendar (between from and to maximum interval can be 3 months)
func (c *CompanyValuation) EarningCalendar(from, to *time.Time) (eList []objects.EarningCalendar, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if from != nil {
		reqParam["from"] = from.Format("2006-01-02")
	}

	if to != nil {
		reqParam["to"] = from.Format("2006-01-02")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationEarningCalendar)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &eList)
	if err != nil {
		return nil, err
	}

	return eList, nil
}

// EarningSurpriseList ...
func (c *CompanyValuation) EarningSurpriseList(symbol string) (eList []objects.EarningSurprise, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationEarningsSurpises, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &eList)
	if err != nil {
		return nil, err
	}

	return eList, nil
}

// EarningCallTranscript - transcript of specific earning
func (c *CompanyValuation) EarningCallTranscript(req objects.RequestEarningCallTranscript) (tList []objects.EarningCallTranscript, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{
			"apikey":  c.apiKey,
			"quarter": fmt.Sprint(req.Quarter),
			"year":    fmt.Sprint(req.Year),
		}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationEarningCallTranscript, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &tList)
	if err != nil {
		return nil, err
	}

	return tList, nil
}

// HistoryEarningCalendar - historical earning calendar
func (c *CompanyValuation) HistoryEarningCalendar(symbol string) (eList []objects.EarningCalendar, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryEarningCalendar, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &eList)
	if err != nil {
		return nil, err
	}

	return eList, nil
}

// IPOCalendar - IPO calendar
func (c *CompanyValuation) IPOCalendar(from, to *time.Time) (ipoList []objects.IPOCalendar, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if from != nil {
		reqParam["from"] = from.Format("2006-01-02")
	}

	if to != nil {
		reqParam["to"] = from.Format("2006-01-02")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationIPOCalendar)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &ipoList)
	if err != nil {
		return nil, err
	}

	return ipoList, nil
}

// SplitCalendar - stock split calendar
func (c *CompanyValuation) SplitCalendar(from, to *time.Time) (sList []objects.SplitCalendar, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if from != nil {
		reqParam["from"] = from.Format("2006-01-02")
	}

	if to != nil {
		reqParam["to"] = from.Format("2006-01-02")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationSplitCalendar)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// DividendCalendar - dividend calendar
func (c *CompanyValuation) DividendCalendar(from, to *time.Time) (dList []objects.DividendCalendar, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if from != nil {
		reqParam["from"] = from.Format("2006-01-02")
	}

	if to != nil {
		reqParam["to"] = from.Format("2006-01-02")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationDividendCalendar)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &dList)
	if err != nil {
		return nil, err
	}

	return dList, nil
}

// InstitutionalHolders - institutional holders
func (c *CompanyValuation) InstitutionalHolders(symbol string) (hList []objects.InstitutionalHolder, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationInstituionalHolder, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &hList)
	if err != nil {
		return nil, err
	}

	return hList, nil
}

// MutualFundHolders - mutual fund holders
func (c *CompanyValuation) MutualFundHolders(symbol string) (hList []objects.MutualFundHolder, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationMutualFundHolder, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &hList)
	if err != nil {
		return nil, err
	}

	return hList, nil
}

// ETFHolders - ETF holders
func (c *CompanyValuation) ETFHolders(symbol string) (hList []objects.ETFHolder, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationETFHolder, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &hList)
	if err != nil {
		return nil, err
	}

	return hList, nil
}

// ETFSectorWeightings - ETF sector weightings
func (c *CompanyValuation) ETFSectorWeightings(symbol string) (sList []objects.ETFSectorWeighting, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationETFSectorWeightings, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// ETFCountryWeightings - ETF country weightings
func (c *CompanyValuation) ETFCountryWeightings(symbol string) (cList []objects.ETFCountryWeighting, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationETFCountryWeightings, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &cList)
	if err != nil {
		return nil, err
	}

	return cList, nil
}

// IncomeStatement - income statement
func (c *CompanyValuation) IncomeStatement(req objects.RequestIncomeStatement) (sList []objects.IncomeStatement, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationIncomeStatement, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// IncomeStatementGrowth - income statement growth
func (c *CompanyValuation) IncomeStatementGrowth(req objects.RequestIncomeStatementGrowth) (sList []objects.IncomeStatementGrowth, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationIncomeStatementGrowth, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// BalanceSheetStatement - balance sheet statement
func (c *CompanyValuation) BalanceSheetStatement(req objects.RequestBalanceSheetStatement) (sList []objects.BalanceSheetStatement, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationBalanceSheetStatement, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// BalanceSheetStatementGrowth - balance sheet statement growth
func (c *CompanyValuation) BalanceSheetStatementGrowth(req objects.RequestBalanceSheetStatementGrowth) (sList []objects.BalanceSheetStatementGrowth, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationBalanceSheetStatementGrowth, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// CashFlowStatement - cash flow statement
func (c *CompanyValuation) CashFlowStatement(req objects.RequestCashFlowStatement) (sList []objects.CashFlowStatement, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationCashFlowStatement, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// CashFlowStatementGrowth - cash flow statement growth
func (c *CompanyValuation) CashFlowStatementGrowth(req objects.RequestCashFlowStatementGrowth) (sList []objects.CashFlowStatementGrowth, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationCashFlowStatementGrowth, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// IncomeStatementAsReported - income statement AS REPORTED
func (c *CompanyValuation) IncomeStatementAsReported(req objects.RequestIncomeStatementAsReported) (sList []objects.IncomeStatementAsReported, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationIncomeStatementAsReported, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// BalanceSheetStatementAsReported - balance sheet statement AS REPORTED
func (c *CompanyValuation) BalanceSheetStatementAsReported(req objects.RequestBalanceSheetStatementAsReported) (sList []objects.BalanceSheetStatementAsReported, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationBalanceSheetStatementAsReported, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// CashFlowStatementAsReported - cash flow statement AS REPORTED
func (c *CompanyValuation) CashFlowStatementAsReported(req objects.RequestCashFlowStatementAsReported) (sList []objects.CashFlowStatementAsReported, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationCashFlowStatementAsReported, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// FullFinancialStatementAsReported - full financial statement AS REPORTED
func (c *CompanyValuation) FullFinancialStatementAsReported(req objects.RequestFullFinancialStatementAsReported) (sList []objects.FullFinancialStatementAsReported, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationFinancialStatementFullAsReported, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// FinancialRatios - financial ratios
func (c *CompanyValuation) FinancialRatios(req objects.RequestFinancialRatios) (rList []objects.FinancialRatios, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationFinancialRatios, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// FinancialRatiosTTM - financial ratios TTM
func (c *CompanyValuation) FinancialRatiosTTM(symbol string) (rList []objects.FinancialRatiosTTM, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationFinancialRatiosTTM, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// KeyMetrics - key metrics
func (c *CompanyValuation) KeyMetrics(req objects.RequestKeyMetrics) (mList []objects.KeyMetrics, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationKeyMetrics, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &mList)
	if err != nil {
		return nil, err
	}

	return mList, nil
}

// KeyMetricsTTM - key metrics ttm
func (c *CompanyValuation) KeyMetricsTTM(symbol string) (mList []objects.KeyMetricsTTM, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationKeyMetricsTTM, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &mList)
	if err != nil {
		return nil, err
	}

	return mList, nil
}

// EnterpriseValue - enterprise value
func (c *CompanyValuation) EnterpriseValue(req objects.RequestEnterpriseValue) (vList []objects.EnterpriseValue, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationEnterpriseValues, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// FinancialStatementsGrowth - financial statements growth
func (c *CompanyValuation) FinancialStatementsGrowth(req objects.RequestFinancialStatementsGrowth) (vList []objects.FinancialStatementsGrowth, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationFinancialGrowth, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// DiscountedCashFlow - discounted cash flow value
func (c *CompanyValuation) DiscountedCashFlow(symbol string) (vList []objects.DiscountedCashFlow, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationDiscountedCashFlow, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// DailyDiscountedCashFlow - daily historical DCF
func (c *CompanyValuation) DailyDiscountedCashFlow(req objects.RequestDailyDiscountedCashFlow) (vList []objects.DailyDiscountedCashFlow, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryDailyDiscountedCashFlow, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// HistoryDiscountedCashFlow - history DCF
func (c *CompanyValuation) HistoryDiscountedCashFlow(req objects.RequestHistoryDiscountedCashFlow) (vList []objects.HistoryDiscountedCashFlow, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryDiscountedCashFlow, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// Rating - get rating by symbol
func (c *CompanyValuation) Rating(symbol string) (rList []objects.Rating, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationRating, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// DailyHistoryRating - daily historical rating
func (c *CompanyValuation) DailyHistoryRating(req objects.RequestRating) (rList []objects.Rating, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryRating, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// MarketCapitalization - market capitalization
func (c *CompanyValuation) MarketCapitalization(symbol string) (rList []objects.MarketCapitalization, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationMarketCapitalization, symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// DailyHistoryMarketCapitalization - daily historical market capitalization
func (c *CompanyValuation) DailyHistoryMarketCapitalization(req objects.RequestMarketCapitalization) (rList []objects.MarketCapitalization, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryMarketCapitalization, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// StockScreener - stock screener
func (c *CompanyValuation) StockScreener(req objects.RequestStockScreener) (sList []objects.StockScreener, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if len(req.Exchange) != 0 {
		reqParam["exchange"] = strings.Join(req.Exchange, ",")
	}

	if req.Industry != nil {
		reqParam["industry"] = string(*req.Industry)
	}

	if req.Sector != nil {
		reqParam["sector"] = string(*req.Sector)
	}

	if req.MarketCapMoreThan != nil {
		reqParam["marketCapMoreThan"] = fmt.Sprint(*req.MarketCapMoreThan)
	}

	if req.MarketCapLowerThan != nil {
		reqParam["marketCapLowerThan"] = fmt.Sprint(*req.MarketCapLowerThan)
	}

	if req.VolumeMoreThan != nil {
		reqParam["volumeMoreThan"] = fmt.Sprint(*req.VolumeMoreThan)
	}

	if req.VolumeLowerThan != nil {
		reqParam["volumeLowerThan"] = fmt.Sprint(*req.VolumeLowerThan)
	}

	if req.BetaMoreThan != nil {
		reqParam["betaMoreThan"] = fmt.Sprint(*req.BetaMoreThan)
	}

	if req.BetaLowerThan != nil {
		reqParam["betaLowerThan"] = fmt.Sprint(*req.BetaLowerThan)
	}

	if req.DividendMoreThan != nil {
		reqParam["dividendMoreThan"] = fmt.Sprint(*req.DividendMoreThan)
	}

	if req.DividendLowerThan != nil {
		reqParam["dividendLowerThan"] = fmt.Sprint(*req.DividendLowerThan)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationStockScreener)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &sList)
	if err != nil {
		return nil, err
	}

	return sList, nil
}

// DelstedCompanies - delsted companies
func (c *CompanyValuation) DelstedCompanies(limit int64) (cList []objects.DelstedCompany, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{
			"apikey": c.apiKey,
			"limit":  fmt.Sprint(limit),
		}).
		Get(c.url + urlAPICompanyValuationDelistedCompanyList)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &cList)
	if err != nil {
		return nil, err
	}

	return cList, nil
}

// StockNews - stock news
func (c *CompanyValuation) StockNews(req objects.RequestStockNews) (vList []objects.StockNews, err error) {
	reqParam := map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}
	if len(req.SymbolList) != 0 {
		reqParam["tickers"] = strings.Join(req.SymbolList, ",")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationStockNews)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// AnalystEstimates - analyst estimates of a stock (Annual || Quarter)
func (c *CompanyValuation) AnalystEstimates(req objects.RequestAnalystEstimates) (vList []objects.AnalystEstimates, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if req.Period != objects.CompanyValuationPeriodAnnual {
		reqParam["period"] = string(objects.CompanyValuationPeriodQuarter)
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationAnalystEstimates, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &vList)
	if err != nil {
		return nil, err
	}

	return vList, nil
}

// Grade - stock grade from analysts
func (c *CompanyValuation) Grade(req objects.RequestGrade) (gList []objects.Grade, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationGrade, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &gList)
	if err != nil {
		return nil, err
	}

	return gList, nil
}

// AnalystStockRecommendations - monthly stock analyst ratings
func (c *CompanyValuation) AnalystStockRecommendations(req objects.RequestAnalystStockRecommendations) (rList []objects.AnalystStockRecommendations, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationAnalystStockRecommendations, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &rList)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

// PressReleases - stock press releases
func (c *CompanyValuation) PressReleases(req objects.RequestPressReleases) (prList []objects.PressReleases, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "limit": fmt.Sprint(req.Limit)}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationPressReleases, req.Symbol))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &prList)
	if err != nil {
		return nil, err
	}

	return prList, nil
}

// FinancialStatementList - List of symbols that have financial statements
func (c *CompanyValuation) FinancialStatementList() (fsList []string, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + urlAPICompanyValuationFinancialStatementsList)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &fsList)
	if err != nil {
		return nil, err
	}

	return fsList, nil
}

// EconomicCalendar - Economic Calendar for time period
func (c *CompanyValuation) EconomicCalendar(req objects.RequestEconomicCalendar) (eList []objects.EconomicCalendar, err error) {
	reqParam := map[string]string{"apikey": c.apiKey}
	if req.From != nil {
		reqParam["from"] = req.From.Format("2006-01-02")
	}

	if req.To != nil {
		reqParam["to"] = req.To.Format("2006-01-02")
	}

	data, err := c.Client.R().
		SetQueryParams(reqParam).
		Get(c.url + urlAPICompanyValuationEconomicCalendar)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &eList)
	if err != nil {
		return nil, err
	}

	return eList, nil
}

// HistoryEconomicCalendar - Economic calendar event list
func (c *CompanyValuation) HistoryEconomicCalendar(req objects.RequestHistoryEconomicCalendar) (hList []objects.HistoryEconomicCalendar, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey, "country": req.Country}).
		Get(c.url + fmt.Sprintf(urlAPICompanyValuationHistoryEconomicCalendar, req.Event))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &hList)
	if err != nil {
		return nil, err
	}

	return hList, nil
}

// EconomicCalendarEventList - Example of historical consumer sentiment in U.S. (take event name and country from event list endpoint)
func (c *CompanyValuation) EconomicCalendarEventList() (eList []objects.EconomicCalendarEventList, err error) {
	data, err := c.Client.R().
		SetQueryParams(map[string]string{"apikey": c.apiKey}).
		Get(c.url + urlAPICompanyValuationEconomicCalendarEventList)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data.Body(), &eList)
	if err != nil {
		return nil, err
	}

	return eList, nil
}
