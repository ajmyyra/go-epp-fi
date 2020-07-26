package epp

import (
	"encoding/xml"
	"time"
)

type APIPoll struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Poll struct {
			Op   string `xml:"op,attr"`
			MsgID string `xml:"msgID,attr,omitempty"`
		} `xml:"poll"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIPollResponse struct {
	XMLName  xml.Name `xml:"epp"`
	Obj      string   `xml:"obj,attr"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result Result `xml:"result"`
		MsgQ PollMessage `xml:"msgQ"`
		ResData struct {
			TrnData struct {
				Name string `xml:"name"`
			} `xml:"trnData"`
		} `xml:"resData"`
		TrID struct {
			ClTRID string `xml:"clTRID"`
			SvTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}

type PollMessage struct {
	Count int    `xml:"count,attr" json:"count"`
	ID    string `xml:"id,attr" json:"id"`
	RawQDate string `xml:"qDate" json:"-"`
	QDate time.Time `json:"date"`
	Msg   string `xml:"msg" json:"msg"`
	Name  string `json:"name"`
}