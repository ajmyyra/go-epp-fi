package registry

import (
	"crypto/tls"
	"fmt"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/gemalto/flume"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	registryServer string
	tlsConfig      tls.Config
	credentials    Credentials

	conn           net.Conn
	readTimeout    time.Duration
	writeTimeout   time.Duration

	log            flume.Logger

	Greeting       epp.Greeting
	LoggedIn       bool
}

type Credentials struct {
	username string
	password string
}

func NewRegistryClient(username, password, serverHost string, serverPort int, clientKey, clientCert []byte) (*Client, error) {
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	registry := fmt.Sprintf("%s:%d", serverHost, serverPort)

	client := Client{
		registryServer: registry,
		readTimeout:    time.Duration(60) * time.Second,
		writeTimeout:   time.Duration(60) * time.Second,
		conn:           nil,
	}
	client.credentials = Credentials{
		username: username,
		password: password,
	}
	client.tlsConfig = tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	loggingConfig := "{\"level\":\"INF\"}"
	if err = flume.ConfigString(loggingConfig); err != nil {
		return nil, err
	}
	if err = flume.ConfigFromEnv(); err != nil {
		return nil, err
	}
	client.log = flume.New("FI EPP")


	// For request ID generation
	rand.Seed(time.Now().UnixNano())

	return &client, nil
}

func (s *Client) SetReadTimeout(seconds int) error {
	if seconds <= 0 {
		return errors.New("Read timeout must be a positive integer.")
	}

	s.readTimeout = time.Duration(seconds) * time.Second

	return nil
}

func (s *Client) SetWriteTimeout(seconds int) error {
	if seconds <= 0 {
		return errors.New("Write timeout must be a positive integer.")
	}

	s.writeTimeout = time.Duration(seconds) * time.Second

	return nil
}
