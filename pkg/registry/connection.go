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
	"net"
	"time"
)

type Client struct {
	RegistryServer string
	TLSConfig tls.Config
	ReadTimeout time.Duration
	WriteTimeout time.Duration

	Conn net.Conn
	Greeting epp.Greeting
	LoggedIn bool
}

func NewRegistryClient(username, password, serverHost string, serverPort int, clientKey, clientCert []byte) (*Client, error) {
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	registry := fmt.Sprintf("%s:%d", serverHost, serverPort)

	client := Client{
		RegistryServer: registry,
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(20) * time.Second,
	}
	client.TLSConfig = tls.Config{
		Certificates: []tls.Certificate{cert},
	}

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
	fmt.Printf("Bytes: %s\n", string(greet)) // DEBUG, remove

	var apigreeting epp.APIGreeting
	if err = xml.Unmarshal(greet, &apigreeting); err != nil {
		return err
	}

	s.Greeting = apigreeting.Greeting
	fmt.Printf("%+v\n", s.Greeting) // TODO logger and to debug

	if s.Greeting.SvcMenu.Version == "1" {
		return errors.New("")
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

func (s *Client) Close() error {
	if err := s.Conn.Close(); err != nil {
		return err
	}

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
