package icici_breeze_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/algotuners-v2/icici-breeze-sdk/breeze_models"
	"github.com/algotuners-v2/icici-breeze-sdk/breeze_models/enums"
	"github.com/algotuners-v2/icici-breeze-sdk/utils"
	"net/http"
	netUrl "net/url"
	"strconv"
	"time"
)

type Breeze struct {
	baseUrl      string
	sessionId    string
	sessionToken string
	userId       string
	password     string
	apiKey       string
	totpCode     string
	environment  string
	pathToSdkDir string
}

type BreezeHistoricalDataInput struct {
	StockCode   string
	ExchCode    string
	FromDate    time.Time
	ToDate      time.Time
	Interval    string
	ProductType string
	ExpiryDate  time.Time
	Right       string
	StrikePrice int
}

const (
	baseBreezeUrl          = "https://breezeapi.icicidirect.com"
	apiBreezeUrl           = "https://api.icicidirect.com/breezeapi/api/v1"
	breezeTimeStringLayout = "2006-01-02T15:04:05.000Z"
	apiLimitOnDataPoints   = 1000
	Mac                    = "mac"
	Linux                  = "linux"
)

var (
	BreezeClient = &Breeze{}
)

func (b *Breeze) Init(userId string, password string, apiKey string, totpCode string, environment string) {
	b.userId = userId
	b.password = password
	b.apiKey = apiKey
	b.totpCode = totpCode
	b.baseUrl = baseBreezeUrl
	b.environment = environment
	b.generateSessionToken(userId, password, apiKey, totpCode)
}

func (b *Breeze) getTimeFormatString(t time.Time) string {
	timeString := t.Format(breezeTimeStringLayout)
	return timeString
}

func (b *Breeze) getMaxNumberOfDataPointsForTimeframe() int {
	return 1000
}

func (b *Breeze) getCurrentDayMarketStartAndEndTime(currentTime time.Time) (time.Time, time.Time) {
	weekDay := currentTime.Weekday()
	if weekDay == time.Sunday || weekDay == time.Saturday {
		currentTime = currentTime.AddDate(0, 0, 1)
		return b.getCurrentDayMarketStartAndEndTime(currentTime)
	}
	marketStartTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 15, 0, 0, currentTime.Location())
	marketEndTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 30, 0, 0, currentTime.Location())

	return marketStartTime, marketEndTime
}

func (b *Breeze) getCurrentDayOrNextMarketStartAndEndTime(currentTime time.Time) (time.Time, time.Time) {
	weekDay := currentTime.Weekday()
	marketStartTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 15, 0, 0, currentTime.Location())
	marketEndTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 30, 0, 0, currentTime.Location())

	if weekDay == time.Sunday || weekDay == time.Saturday {
		currentTime = currentTime.AddDate(0, 0, 1)
		return b.getCurrentDayMarketStartAndEndTime(currentTime)
	}

	if currentTime.Equal(marketEndTime) || currentTime.After(marketEndTime) {
		nextDay := currentTime.AddDate(0, 0, 1)
		nextDayMarketStartTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 15, 0, 0, currentTime.Location())
		return b.getCurrentDayOrNextMarketStartAndEndTime(nextDayMarketStartTime)
	}

	return marketStartTime, marketEndTime
}

func (b *Breeze) findMinTime(time1, time2 time.Time) time.Time {
	if time1.Before(time2) {
		return time1
	}
	return time2
}

func (b *Breeze) findMaxTime(time1, time2 time.Time) time.Time {
	if time1.After(time2) {
		return time1
	}
	return time2
}

func (b *Breeze) getEndTimeUnderApiLimitForTimestamp(startTime time.Time, timeframe string) time.Time {
	_, marketEndTime := b.getCurrentDayMarketStartAndEndTime(startTime)
	switch timeframe {
	case enums.OneSecond:
		numberOfMinutes := (apiLimitOnDataPoints / 60) - 1
		calculatedEndTime := startTime.Add(time.Duration(numberOfMinutes) * time.Minute)
		if calculatedEndTime.Before(marketEndTime) {
			return calculatedEndTime
		}
		return marketEndTime.Add(time.Duration(1) * time.Hour)
	case enums.OneMinute:
		numberOfHours := 2 * 24
		calculatedEndTime := startTime.Add(time.Duration(numberOfHours) * time.Hour)
		return calculatedEndTime
	case enums.FiveMinutes:
		numberOfHours := 13 * 24
		calculatedEndTime := startTime.Add(time.Duration(numberOfHours) * time.Hour)
		return calculatedEndTime
	case enums.Day:
		numberOfHours := 900 * 24
		calculatedEndTime := startTime.Add(time.Duration(numberOfHours) * time.Hour)
		return calculatedEndTime
	}
	return startTime.Add(time.Duration(999999) * time.Hour)
}

func (b *Breeze) GetHistoricalData(input BreezeHistoricalDataInput) []breeze_models.CandleData {
	url := fmt.Sprintf("%s/api/v2/historicalcharts?stock_code=%s&exch_code=%s&from_date=%s&to_date=%s&interval=%s&product_type=%s&expiry_date=%s&right=%s&strike_price=%s", b.baseUrl, input.StockCode, input.ExchCode, b.getTimeFormatString(input.FromDate), b.getTimeFormatString(input.ToDate), input.Interval, input.ProductType, b.getTimeFormatString(input.ExpiryDate), input.Right, strconv.Itoa(input.StrikePrice))
	var headers http.Header = map[string][]string{}
	queryParams := netUrl.Values{
		"stock_code":   {input.StockCode},
		"exch_code":    {input.ExchCode},
		"from_date":    {b.getTimeFormatString(input.FromDate)},
		"to_date":      {b.getTimeFormatString(input.ToDate)},
		"interval":     {input.Interval},
		"product_type": {input.ProductType},
		"expiry_date":  {b.getTimeFormatString(input.ExpiryDate)},
		"right":        {input.Right},
		"strike_price": {strconv.Itoa(input.StrikePrice)},
	}
	headers.Add("Content-Type", "application/json")
	headers.Add("apikey", b.apiKey)
	headers.Add("X-SessionToken", b.sessionToken)
	httpClient := utils.GenerateHttpClient(nil, false)
	res, err := httpClient.Do(http.MethodGet, url, queryParams, headers)
	utils.Log(err)
	var targetStruct breeze_models.HistoricalDataResponse
	err = json.Unmarshal(res.Body, &targetStruct)
	if targetStruct.Status != 200 {
		utils.Log(errors.New(targetStruct.Error))
	}
	utils.Log(err)
	return targetStruct.Success
}

func (b *Breeze) getDurationTimeframe(timeframe string) time.Duration {
	switch timeframe {
	case enums.OneSecond:
		return time.Second
	case enums.OneMinute:
		return time.Minute
	case enums.FiveMinutes:
		return time.Minute * time.Duration(5)
	case enums.Day:
		return time.Hour * time.Duration(24)
	}
	utils.Log(errors.New("Invalid timeframe"))
	return time.Hour * 24 * 30 * 12 * 100
}

func (b *Breeze) isTimeAfterMarketHours(timestamp time.Time) bool {
	if timestamp.Hour() > 15 {
		return true
	}
	if timestamp.Hour() == 15 && timestamp.Minute() > 30 {
		return true
	}
	return false
}

func (b *Breeze) GetHistoricalDataV2(input BreezeHistoricalDataInput) []breeze_models.CandleData {
	startTime := input.FromDate
	endTime := input.ToDate
	interval := input.Interval
	currentTime := startTime
	candlesData := []breeze_models.CandleData{}
	for currentTime.Before(endTime) {
		currentEndTime := b.getEndTimeUnderApiLimitForTimestamp(currentTime, interval)
		currentEndTime = b.findMinTime(currentEndTime, endTime)
		input.FromDate = currentTime
		input.ToDate = currentEndTime
		candles := b.GetHistoricalData(input)
		candlesData = append(candlesData, candles...)
		currentTime = currentEndTime
		currentTime = currentTime.Add(b.getDurationTimeframe(interval))
		currentTimeWeekDay := currentTime.Weekday()
		if b.isTimeAfterMarketHours(currentTime) || currentTimeWeekDay == time.Saturday || currentTimeWeekDay == time.Sunday {
			currentTime, _ = b.getCurrentDayOrNextMarketStartAndEndTime(currentTime)
		}
		time.Sleep(time.Second)
	}
	return candlesData
}
