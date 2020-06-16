package registry

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/xml"
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

func (s *Client) Connect() error {
	dialConn, err := tls.Dial("tcp", s.registryServer, &s.tlsConfig)
	if err != nil {
		return err
	}
	s.conn = dialConn

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

	if s.readTimeout > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	}

	err := binary.Read(s.conn, binary.BigEndian, &rawResponse)
	if err != nil {
		return nil, err
	}

	rawResponse -= 4
	if rawResponse < 0 {
		return nil, io.ErrUnexpectedEOF
	}

	bytesResponse, err := readStreamToBytes(s.conn, rawResponse)
	if err != nil {
		return nil, err
	}

	return bytesResponse, nil
}

func (s *Client) Write(payload []byte) error {
	payload = []byte(xml.Header + string(payload))

	sendBytesLength := uint32(4 + len(payload))

	if s.writeTimeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	}

	err := binary.Write(s.conn, binary.BigEndian, sendBytesLength)
	if err != nil {
		return err
	}
	if _, err = s.conn.Write(payload); err != nil {
		return err
		// TODO log first param (amount of bytes written) if error
	}

	return nil
}

func (s *Client) Send(payload []byte) ([]byte, error) {
	s.log.Debug("Sending message:\n" + string(payload))
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

	s.log.Debug("Received response:\n" + string(apiResp))

	return apiResp, nil
}

func (s *Client) Close() error {
	if err := s.conn.Close(); err != nil {
		return err
	}

	s.conn = nil
	return nil
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

func createRequestID(length int) string {
	reqID := make([]byte, length)
	for i := range reqID {
		reqID[i] = reqIDChars[rand.Intn(len(reqIDChars))]
	}
	return string(reqID)
}

func parseDate(rawDate string) (time.Time, error) {
	emptyDateFormat := "0001-01-01T00:00:00"
	greetingDateFormat := time.RFC3339Nano
	pollDateFormat := "2006-01-02T15:04:05"
	domainDateFormat := "2006-01-02T15:04:05.000"
	renewalDateFormat := "2006-01-03T15:04:05.0Z"

	if rawDate == emptyDateFormat {
		return time.Time{}, nil
	}

	if date, err := time.Parse(greetingDateFormat, rawDate); err == nil {
		return date, nil
	}
	if date, err := time.Parse(domainDateFormat, rawDate); err == nil {
		return date, nil
	}
	if date, err := time.Parse(pollDateFormat, rawDate); err == nil {
		return date, nil
	}
	if date, err := time.Parse(renewalDateFormat, rawDate); err == nil {
		return date, nil
	}

	return time.Time{}, errors.New("Unrecognised date format: " + rawDate)
}

func (s *Client) logAPIConnectionError(err error, args ...string) {
	s.log.Error("API connection failed when making a request", "error", err, args)
}

