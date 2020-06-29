package epp

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"net"
	"time"
)

type APIHostCheck struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Check struct {
			HostCheck struct {
				Xmlns  string `xml:"xmlns:host,attr"`
				Name   []string `xml:"host:name"`
			} `xml:"host:check"`
		} `xml:"check"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIHostInfo struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Info struct {
			HostInfo struct {
				Xmlns string `xml:"xmlns:host,attr"`
				Name string `xml:"host:name"`
			} `xml:"host:info"`
		} `xml:"info"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIHostInfoResponse struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`
		ResData struct {
			HostInfo HostInfoResp `xml:"infData"`
		} `xml:"resData"`
		TrID Transaction `xml:"trID"`
	} `xml:"response"`
}

type APIHostCreation struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Create struct {
			HostCreate struct {
				Xmlns string `xml:"xmlns:host,attr"`
				Hostname string `xml:"host:name"`
				Addr []HostIPAddress `xml:"host:addr"`
			} `xml:"host:create"`
		} `xml:"create"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIHostUpdate struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Update struct {
			HostUpdate struct {
				Xmlns string `xml:"xmlns:host,attr"`
				Hostname string `xml:"host:name"`
				Add  struct {
					Addr []HostIPAddress `xml:"host:addr,omitempty"`
				} `xml:"host:add"`
				Rem struct {
					Addr []HostIPAddress `xml:"host:addr,omitempty"`
				} `xml:"host:rem"`
			} `xml:"host:update"`
		} `xml:"update"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIHostDeletion struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Delete struct {
			HostDelete struct {
				Xmlns string `xml:"xmlns:host,attr"`
				Hostname string `xml:"host:name"`
			} `xml:"host:delete"`
		} `xml:"delete"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type HostInfoResp struct {
	Xmlns     string `xml:"host,attr"`
	Hostname  string `xml:"name"`
	Addr      []struct {
		IP       string `xml:",chardata"`
		Family   string `xml:"ip,attr"`
	} `xml:"addr"`
	ClID      string `xml:"clID"`
	CrID      string `xml:"crID"`
	RawCrDate string `xml:"crDate"`
	CrDate    time.Time
	RawUpDate string `xml:"upDate"`
	UpDate    time.Time
}

type HostIPAddress struct {
	IP     string `xml:",chardata"`
	Family string `xml:"ip,attr"`
}

func FormatHostIP(rawIp string) (HostIPAddress, error) {
	ip := net.ParseIP(rawIp)
	if ip.To4() != nil {
		return HostIPAddress{
			IP:     rawIp,
			Family: "v4",
		}, nil
	} else if ip.To16() != nil {
		return HostIPAddress{
			IP:     rawIp,
			Family: "v6",
		}, nil
	} else {
		return HostIPAddress{}, errors.New("Unrecognised IP address format: " + rawIp)
	}
}