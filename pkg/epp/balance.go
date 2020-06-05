package epp

import "encoding/xml"

type APIBalance struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Check struct {
			Balance string `xml:"balance"`
		} `xml:"check"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}
