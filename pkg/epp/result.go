package epp

import "encoding/xml"

type APIResult struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result Result `xml:"result"`
		TrID   Transaction `xml:"trID"`
	} `xml:"response"`
}

type Result struct {
	Code int `xml:"code,attr"`
	Msg  string `xml:"msg"`
}

type Transaction struct {
	ClTRID string `xml:"clTRID"`
	SvTRID string `xml:"svTRID"`
}
