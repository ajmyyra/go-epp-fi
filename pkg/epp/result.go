package epp

import (
	"encoding/xml"
	"time"
)

type APIResult struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result Result `xml:"result"`
		ResData ResData `xml:"resData"`
		TrID   Transaction `xml:"trID"`
	} `xml:"response"`
}

type Result struct {
	Code int `xml:"code,attr"`
	Msg  string `xml:"msg"`
}

type ResData struct {
	BalanceAmount int    `xml:"balanceamount"`
	Timestamp     string `xml:"timestamp"`
	ChkData       struct {
		Cd []ItemCheck `xml:"cd"`
	} `xml:"chkData"`
	CreateData  CreateData `xml:"creData"`
	RenewalData RenewalData `xml:"renData"`
	TransferData TransferData `xml:"trnData"`
}

type CreateData struct {
	ID        string `xml:"id"`
	Name      string `xml:"name"`
	RawCrDate string `xml:"crDate"`
	CrDate    time.Time
	RawExDate string `xml:"exDate"`
	ExDate    time.Time
}

type RenewalData struct {
	Name       string `xml:"name"`
	RawExpDate string `xml:"exDate"`
	ExpireDate time.Time
}

type TransferData struct {
	Xmlns     string `xml:"obj,attr"`
	Name      string `xml:"name"`
	TrStatus  string `xml:"trStatus"`
	ReID      string `xml:"reID"`
	ReRawDate string `xml:"reDate"`
	ReDate    time.Time
	AcID      string `xml:"acID"`
}

type ItemCheck struct {
	ContactId struct {
		Name  string `xml:",chardata"`
		Avail int    `xml:"avail,attr"`
	} `xml:"id"`
	Name struct {
		Name  string `xml:",chardata"`
		Avail int    `xml:"avail,attr"`
	} `xml:"name"`
	Reason      string `xml:"reason"`
	IsAvailable bool
}

type Transaction struct {
	ClTRID string `xml:"clTRID"`
	SvTRID string `xml:"svTRID"`
}
