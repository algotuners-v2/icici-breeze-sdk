package breeze_models

import (
	"reflect"
	"strconv"
	"time"
)

type NseScrip struct {
	Token                     string    `csv:"Token"`
	ShortName                 string    `csv:"ShortName"`
	Series                    string    `csv:"Series"`
	CompanyName               string    `csv:"CompanyName"`
	TickSize                  float64   `csv:"ticksize"`
	LotSize                   int       `csv:"Lotsize"`
	DateOfListing             time.Time `csv:"DateOfListing"`
	DateOfDeListing           time.Time `csv:"DateOfDeListing"`
	IssuePrice                float64   `csv:"IssuePrice"`
	FaceValue                 float64   `csv:"FaceValue"`
	ISINCode                  string    `csv:"ISINCode"`
	Weeks52High               float64   `csv:"52WeeksHigh"`
	Weeks52Low                float64   `csv:"52WeeksLow"`
	LifeTimeHigh              float64   `csv:"LifeTimeHigh"`
	LifeTimeLow               float64   `csv:"LifeTimeLow"`
	HighDate                  time.Time `csv:"HighDate"`
	LowDate                   time.Time `csv:"LowDate"`
	Symbol                    string    `csv:"Symbol"`
	InstrumentType            string    `csv:"InstrumentType"`
	PermittedToTrade          string    `csv:"PermittedToTrade"`
	IssueCapital              float64   `csv:"IssueCapital"`
	WarningPercent            float64   `csv:"WarningPercent"`
	FreezePercent             float64   `csv:"FreezePercent"`
	CreditRating              string    `csv:"CreditRating"`
	IssueRate                 float64   `csv:"IssueRate"`
	IssueStartDate            time.Time `csv:"IssueStartDate"`
	InterestPaymentDate       time.Time `csv:"InterestPaymentDate"`
	IssueMaturityDate         time.Time `csv:"IssueMaturityDate"`
	BoardLotQty               int       `csv:"BoardLotQty"`
	Name                      string    `csv:"Name"`
	ListingDate               time.Time `csv:"ListingDate"`
	ExpulsionDate             time.Time `csv:"ExpulsionDate"`
	ReAdmissionDate           time.Time `csv:"ReAdmissionDate"`
	RecordDate                time.Time `csv:"RecordDate"`
	ExpiryDate                time.Time `csv:"ExpiryDate"`
	NoDeliveryStartDate       time.Time `csv:"NoDeliveryStartDate"`
	NoDeliveryEndDate         time.Time `csv:"NoDeliveryEndDate"`
	MFill                     int       `csv:"MFill"`
	AON                       int       `csv:"AON"`
	ParticipantInMarketIndex  string    `csv:"ParticipantInMarketIndex"`
	BookClsStartDate          time.Time `csv:"BookClsStartDate"`
	BookClsEndDate            time.Time `csv:"BookClsEndDate"`
	EGM                       int       `csv:"EGM"`
	AGM                       int       `csv:"AGM"`
	Interest                  float64   `csv:"Interest"`
	Bonus                     int       `csv:"Bonus"`
	Rights                    int       `csv:"Rights"`
	Dividends                 float64   `csv:"Dividends"`
	LocalUpdateDateTime       time.Time `csv:"LocalUpdateDateTime"`
	DeleteFlag                int       `csv:"DeleteFlag"`
	Remarks                   string    `csv:"Remarks"`
	NormalMarketStatus        string    `csv:"NormalMarketStatus"`
	OddLotMarketStatus        string    `csv:"OddLotMarketStatus"`
	SpotMarketStatus          string    `csv:"SpotMarketStatus"`
	AuctionMarketStatus       string    `csv:"AuctionMarketStatus"`
	NormalMarketEligibility   string    `csv:"NormalMarketEligibility"`
	OddLotMarketEligibility   string    `csv:"OddLotlMarketEligibility"`
	SpotMarketEligibility     string    `csv:"SpotMarketEligibility"`
	AuctionlMarketEligibility string    `csv:"AuctionlMarketEligibility"`
	MarginPercentage          float64   `csv:"MarginPercentage"`
	ExchangeCode              string    `csv:"ExchangeCode"`
}

func (ns *NseScrip) MapCSVToStruct(record []string, target interface{}, header []string) error {
	//r := csv.NewReader(strings.NewReader(strings.Join(header, ",")))
	//r.Comma = ','
	//r.LazyQuotes = true
	//// Read the CSV record
	//recordFields, err := r.Read()
	//if err != nil {
	//	return err
	//}

	// Create a map to store the mapping of struct field names to their index in the CSV record
	fieldMap := make(map[string]int)
	for i, field := range header {
		fieldMap[field[2:len(field)-1]] = i
	}

	// Use reflection to set struct fields based on the CSV record
	structValue := reflect.ValueOf(target).Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		csvTag := field.Tag.Get("csv")
		csvIndex, ok := fieldMap[csvTag]
		if ok {
			fieldValue := structValue.Field(i)
			fieldType := fieldValue.Type()

			// Convert the CSV field value to the appropriate struct field type
			switch fieldType.Kind() {
			case reflect.Int:
				intValue, err := strconv.Atoi(record[csvIndex])
				if err != nil {
					fieldValue.SetInt(int64(0))
				} else {
					fieldValue.SetInt(int64(intValue))
				}
			case reflect.Float64:
				floatValue, err := strconv.ParseFloat(record[csvIndex], 64)
				if err != nil {
					fieldValue.SetFloat(0)
				} else {
					fieldValue.SetFloat(floatValue)
				}
			case reflect.String:
				fieldValue.SetString(record[csvIndex])
			case reflect.Struct:
				timeValue, err := time.Parse("2006/01/02 15:04:05", record[csvIndex])
				if err != nil {
					fieldValue.Set(reflect.ValueOf(time.Time{}))
				} else {
					fieldValue.Set(reflect.ValueOf(timeValue))
				}
			}
		}
	}
	return nil
}
