package main

import (
	"fmt"
	icici_breeze_sdk "github.com/algotuners-v2/icici-breeze-sdk"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

func main() {
	breezeClient := icici_breeze_sdk.Breeze{}
	breezeClient.Init(
		"8950577400",
		"Lzt2a7TBbcWH",
		"342_232u02y7541998~582e6762HG787",
		"OBIDM6CCMJIXSRCDNQ2US4K2GI",
		icici_breeze_sdk.Linux,
	)
	startTime := time.Date(2024, 1, 8, 9, 15, 0, 0, time.Local)
	endTime := time.Date(2024, 1, 8, 9, 16, 0, 0, time.Local)
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
	////breezeClient.GetHistoricalDataV2(iciciInput)
	fmt.Println(len(data))

	//scrips := icici_breeze_sdk.GetNseFnoScrips()
	//fmt.Println(scrips[1:2])
}

func main1() {
	service, err := selenium.NewChromeDriverService("/Users/mayanksheoran/Desktop/Root/Coding/Github/mayank/algotuners-v2/icici-breeze-sdk/chromedriver-mc", 4444)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		"--headless", // comment out this line to see the browser
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}

	err = driver.Get("https://www.google.com")
	if err != nil {
		return
	}
}
