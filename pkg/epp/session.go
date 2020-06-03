package epp

import "encoding/xml"

type APILogin struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Login Login `xml:"login"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type Login struct {
	ClID    string `xml:"clID"`
	Pw      string `xml:"pw"`
	NewPW   string `xml:"newPW,omitempty"`
	Options struct {
		Version string `xml:"version"`
		Lang    string `xml:"lang"`
	} `xml:"options"`
	Svcs struct {
		ObjURI       []string `xml:"objURI"`
		SvcExtension struct {
			ExtURI []string `xml:"extURI"`
		} `xml:"svcExtension"`
	} `xml:"svcs"`
}

type APILogout struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Logout string `xml:"logout"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}