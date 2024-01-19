package icici_breeze_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/algotuners-v2/icici-breeze-sdk/breeze_models"
	"github.com/algotuners-v2/icici-breeze-sdk/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func (b *Breeze) generateTOTP(secret string) (string, error) {
	otpCode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return otpCode, nil
}

func (b *Breeze) getSessionId(userId string, password string, totpCode string) string {
	dir, _ := os.Getwd()
	chromeDriverName := "chromedriver"
	if b.environment == Mac {
		chromeDriverName = "chromedriver-mc"
	}
	chromePath := path.Join(dir, chromeDriverName)
	service, err := selenium.NewChromeDriverService(chromePath, 4444)
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
		"--headless",
	}})

	chromeDrivers, err := selenium.NewRemote(caps, "")

	utils.Log(err)
	if err == nil {
		fmt.Println("chrome Driver done.")
	}
	defer chromeDrivers.Quit()
	err = chromeDrivers.Get("https://api.icicidirect.com/apiuser/login?api_key=342_232u02y7541998~582e6762HG787")
	utils.Log(err)
	_ = chromeDrivers.SetImplicitWaitTimeout(time.Second * 5)
	userName, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[1]/input")
	utils.Log(err)
	passWord, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[3]/div/input")
	utils.Log(err)
	err = userName.SendKeys(userId)
	utils.Log(err)
	err = passWord.SendKeys(password)
	utils.Log(err)

	// checkbox
	item, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[4]/div/input")
	utils.Log(err)
	err = item.Click()
	utils.Log(err)

	// login
	item, err = chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[5]/input[1]")
	utils.Log(err)
	err = item.Click()
	utils.Log(err)

	pin, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[2]/div/div[2]/div[2]/div[3]/div/div[1]/input[1]")
	otpValue, err := b.generateTOTP(totpCode)
	err = pin.SendKeys(otpValue)
	utils.Log(err)

	item, err = chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[2]/div/div[2]/div[2]/div[4]/input[1]")
	utils.Log(err)
	err = item.Click()
	utils.Log(err)

	time.Sleep(time.Second * 2)
	currentUrl, err := chromeDrivers.CurrentURL()
	utils.Log(err)
	strList := strings.Split(currentUrl, "apisession=")
	sessionId := strList[1]
	return sessionId
}

func (b *Breeze) generateSessionToken(userId string, password string, apiKey string, totpCode string) {
	sessionId := b.getSessionId(userId, password, totpCode)
	customerdetailsUrl := apiBreezeUrl + "/customerdetails"
	jsonBody := strings.NewReader(fmt.Sprintf(`{
		"AppKey": "%s",
		"SessionToken": "%s"
	}`, apiKey, sessionId))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, customerdetailsUrl, jsonBody)
	utils.Log(err)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	utils.Log(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	utils.Log(err)
	cdResponse := breeze_models.CustomerDetailsResponse{}
	err = json.Unmarshal(body, &cdResponse)
	utils.Log(err)
	b.sessionId = sessionId
	b.sessionToken = cdResponse.Success.SessionToken
}
