package epp

import "encoding/xml"

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
		Cd        []ItemCheck `xml:"cd"`
	} `xml:"chkData"`
	CreateData CreateData `xml:"creData"`
}

type CreateData struct {
	ID         string `xml:"id"`
	Name       string `xml:"name"`
	CreateDate string `xml:"crDate"`
	ExpireDate string `xml:"exDate"`
}

type ItemCheck struct {
	ContactId struct {
		Name  string `xml:",chardata"`
		Avail int    `xml:"avail,attr"`
	} `xml:"id"`
	Domain struct {
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
