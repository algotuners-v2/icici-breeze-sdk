package main

import (
	"fmt"
	"github.com/algotuners-v2/icici-breeze-sdk"
	"os"
	"time"
)

func main() {
	breezeClient := icici_breeze_sdk.Breeze{}
	wd, _ := os.Getwd()
	breezeClient.Init(
		"8950577400",
		"Lzt2a7TBbcWH",
		"342_232u02y7541998~582e6762HG787",
		"OBIDM6CCMJIXSRCDNQ2US4K2GI",
		icici_breeze_sdk.Mac,
		wd,
	)
	startTime := time.Date(2024, 1, 8, 9, 15, 0, 0, time.Local)
	endTime := time.Date(2024, 1, 8, 15, 30, 0, 0, time.Local)
	expiryTime := time.Date(2024, 1, 11, 0, 0, 0, 0, time.Local)
	iciciInput := icici_breeze_sdk.BreezeHistoricalDataInput{
		StockCode:   "NIFTY",
		ExchCode:    "NFO",
		FromDate:    startTime,
		ToDate:      endTime,
		Interval:    "1second",
		ProductType: "Options",
		ExpiryDate:  expiryTime,
		Right:       "Call",
		StrikePrice: 21700,
	}
	data := breezeClient.GetHistoricalData(iciciInput)
	//breezeClient.GetHistoricalDataV2(iciciInput)
	fmt.Println(len(data))
}
