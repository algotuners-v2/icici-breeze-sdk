package icici_breeze_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tebeka/selenium"
	"icici-breeze-sdk/breeze_models"
	"icici-breeze-sdk/utils"
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
	caps := selenium.Capabilities{"browserName": "chrome"}
	dir, _ := os.Getwd()
	chromeDriverName := "chromedriver"
	if b.environment == Mac {
		chromeDriverName = "chromedriver-mc"
	}
	chromePath := path.Join(dir, chromeDriverName)
	opts := []selenium.ServiceOption{}
	drivers, err := selenium.NewChromeDriverService(chromePath, 8080, opts...)
	utils.Log(err)
	defer drivers.Stop()
	chromeDrivers, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 8080))
	utils.Log(err)
	defer chromeDrivers.Quit()
	err = chromeDrivers.Get("https://api.icicidirect.com/apiuser/login?api_key=342_232u02y7541998~582e6762HG787")
	utils.Log(err)
	_ = chromeDrivers.SetImplicitWaitTimeout(time.Second * 5)
	userName, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[1]/input")
	utils.Log(err)
	passWord, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[3]/div/input")
	utils.Log(err)
	_ = userName.SendKeys(userId)
	_ = passWord.SendKeys(password)

	// checkbox
	item, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[4]/div/input")
	utils.Log(err)
	_ = item.Click()

	// login
	item, err = chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[1]/div[2]/div/div[5]/input[1]")
	utils.Log(err)
	_ = item.Click()

	pin, err := chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[2]/div/div[2]/div[2]/div[3]/div/div[1]/input[1]")
	otpValue, err := b.generateTOTP(totpCode)
	_ = pin.SendKeys(otpValue)

	item, err = chromeDrivers.FindElement(selenium.ByXPATH, "/html/body/form/div[2]/div/div/div[2]/div/div[2]/div[2]/div[4]/input[1]")
	utils.Log(err)
	_ = item.Click()

	time.Sleep(time.Second * 2)
	currentUrl, _ := chromeDrivers.CurrentURL()
	strList := strings.Split(currentUrl, "apisession=")
	sessionId := strList[1]
	return sessionId
}

func (b *Breeze) generateSessionToken(userId string, password string, apiKey string, totpCode string) {
	sessionId := b.getSessionId(userId, password, totpCode)
	customerdetailsUrl := apiBreezeUrl + "/customerdetails"
	fmt.Println(sessionId)
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
	fmt.Println(cdResponse.Success.SessionToken)
	b.sessionId = sessionId
	b.sessionToken = cdResponse.Success.SessionToken
}
