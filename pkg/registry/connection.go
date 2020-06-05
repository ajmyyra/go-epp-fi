package registry

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"net"
	"time"
)

const APIVersion = "1.0"
const APILanguage = "en"

const reqIDChars = "ABCDEFGHIJKLMNOPQRSTUVXYZW0123456789"
const reqIDLength = 5

type Client struct {
	RegistryServer string
	TLSConfig tls.Config
	Credentials Credentials

	ReadTimeout time.Duration
	WriteTimeout time.Duration

	Conn net.Conn
	Greeting epp.Greeting
	LoggedIn bool
}

type Credentials struct {
	Username string
	Password string
}

func NewRegistryClient(username, password, serverHost string, serverPort int, clientKey, clientCert []byte) (*Client, error) {
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	registry := fmt.Sprintf("%s:%d", serverHost, serverPort)

	client := Client{
		RegistryServer: registry,
		ReadTimeout:    time.Duration(60) * time.Second,
		WriteTimeout:   time.Duration(60) * time.Second,
		Conn: nil,
	}
	client.Credentials = Credentials{
		Username: username,
		Password: password,
	}
	client.TLSConfig = tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// For request ID generation
	rand.Seed(time.Now().UnixNano())

	return &client, nil
}


func (s *Client) Connect() error {
	dialConn, err := tls.Dial("tcp", s.RegistryServer, &s.TLSConfig)
	if err != nil {
		return err
	}
	s.Conn = dialConn

	greet, err := s.Read()
	if err != nil {
		return err
	}

	s.Greeting, err = unmarshalGreeting(greet)
	if err != nil {
		return err
	}

	if s.Greeting.SvcMenu.Version != APIVersion {
		return errors.New("Unexpected version: " + s.Greeting.SvcMenu.Version)
	}

	return nil
}

func (s *Client) Read() ([]byte, error) {
	var rawResponse int32

	if s.ReadTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
	}

	err := binary.Read(s.Conn, binary.BigEndian, &rawResponse)
	if err != nil {
		return nil, err
	}

	rawResponse -= 4
	if rawResponse < 0 {
		return nil, io.ErrUnexpectedEOF
	}

	bytesResponse, err := readStreamToBytes(s.Conn, rawResponse)
	if err != nil {
		return nil, err
	}

	return bytesResponse, nil
}

func (s *Client) Write(payload []byte) error {
	payload = []byte(xml.Header + string(payload))

	sendBytesLength := uint32(4 + len(payload))

	if s.WriteTimeout > 0 {
		s.Conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
	}

	err := binary.Write(s.Conn, binary.BigEndian, sendBytesLength)
	if err != nil {
		return err
	}
	if _, err = s.Conn.Write(payload); err != nil {
		return err
		// TODO log first param (amount of bytes written) if error
	}

	return nil
}

func (s *Client) Send(payload []byte) ([]byte, error) {
	err := s.Write(payload)
	if err != nil {
		// TODO log error
		return nil, err
	}

	time.Sleep(time.Duration(1) * time.Second)

	apiResp, err := s.Read()
	if err != nil {
		// TODO log error
		return nil, err
	}

	return apiResp, nil
}

func (s *Client) Close() error {
	if err := s.Conn.Close(); err != nil {
		return err
	}

	s.Conn = nil
	return nil
}

func (s *Client) Hello() (epp.Greeting, error) {
	if s.Conn != nil {
		hello := epp.APIHello{
			XMLName: xml.Name{},
			Xmlns:   epp.EPPNamespace,
		}
		helloMsg, _ := xml.MarshalIndent(hello, "", "  ")
		apiResp, err := s.Send(helloMsg)
		if err != nil {
			return epp.Greeting{}, err
		}

		greeting, err := unmarshalGreeting(apiResp)
		if err != nil {
			return epp.Greeting{}, err
		}

		return greeting, nil
	}

	return epp.Greeting{}, errors.New("Uninitialized connection, unable to connect to server.")
}

func readStreamToBytes(conn net.Conn, rawResponse int32) ([]byte, error) {
	lr := io.LimitedReader{R: conn, N: int64(rawResponse)}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(&lr); err != nil {
		return nil, err
		// TODO log first param (amount of bytes read) if error
	}
	return buf.Bytes(), nil
}

func unmarshalGreeting(rawGreeting []byte) (epp.Greeting, error) {
	var greeting epp.APIGreeting
	if err := xml.Unmarshal(rawGreeting, &greeting); err != nil {
		return epp.Greeting{}, err
	}

	formattedDate, err := parseDate(greeting.Greeting.RawDate, time.RFC3339Nano)
	if err != nil {
		return epp.Greeting{}, errors.Wrap(err, "Invalid or non-existent date in greeting")
	}
	greeting.Greeting.SvDate = formattedDate

	return greeting.Greeting, nil
}

func createRequestID(length int) string {
	reqID := make([]byte, length)
	for i := range reqID {
		reqID[i] = reqIDChars[rand.Intn(len(reqIDChars))]
	}
	return string(reqID)
}

func parseDate(rawDate string, timeFormat string) (time.Time, error) {
	formattedDate, err := time.Parse(timeFormat, rawDate)
	if err != nil {
		return time.Time{}, err
	}

	return formattedDate, nil
}