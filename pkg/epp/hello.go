package epp

import "encoding/xml"

type APIHello struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Hello   string   `xml:"hello"`
}
