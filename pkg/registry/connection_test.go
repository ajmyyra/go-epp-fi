package registry

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
	"testing"
)

func TestConnectionAndHello(t *testing.T) {
	eppTestServer, err := createEPPTestServer("127.0.0.1", 12001)
	if err != nil {
		t.Fatalf("Error when creating server for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestClient, err := createEPPTestClient("test", "test123", "127.0.0.1", 12001)
	if err != nil {
		t.Fatalf("Error when creating client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(helloReq, greeting, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if _, err = eppTestClient.Hello(); err != nil {
		t.Fatalf("Hello failed: %v\n", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the connection failed: %v\n", err)
	}
}

func createEPPTestClient(user, password, serverHost string, serverPort int) (*Client, error) {
	clientCert, err := ioutil.ReadFile("../../testtmp/testclient.crt")
	if err != nil {
		return nil, errors.Wrap(err, "Have the test certificates been created by running 'make create-test-certs'?")
	}
	clientKey, err := ioutil.ReadFile("../../testtmp/testclient.key")
	if err != nil {
		return nil, errors.Wrap(err, "Have the test certificates been created by running 'make create-test-certs'?")
	}
	caCert, err := ioutil.ReadFile("../../testtmp/rootCA.crt")
	if err != nil {
		return nil, errors.Wrap(err, "Have the test certificates been created by running 'make create-test-certs'?")
	}

	client, err := NewRegistryClient(user, password, serverHost, serverPort, clientKey, clientCert)
	if err != nil {
		return nil, errors.Wrap(err, "Problem creating a new test client.")
	}
	if err = client.SetCACertificates(caCert); err != nil {
		return nil, errors.Wrap(err, "Problem setting CA certificates for the new client.")
	}

	return client, nil
}

type EPPTestServer struct {
	listener    net.Listener

	expectedReq chan []byte
	successResp chan []byte
	errorResp   chan []byte

}

func createEPPTestServer(serverHost string, serverPort int) (EPPTestServer, error) {
	cert, caCert, err := loadTestServerCerts()
	if err != nil {
		return EPPTestServer{}, err
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: caCert,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: caCert,
	}

	listenAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)
	listener, err := tls.Listen("tcp", listenAddr, &config)
	if err != nil {
		return EPPTestServer{}, err
	}

	// As requests and responses can be rather long, we play it safe
	// by having really big channel buffers. Not that nice & clean
	// with production code, but sufficient for testing.
	eppTest := EPPTestServer{
		listener: listener,
		expectedReq: make(chan []byte, 10000),
		successResp: make(chan []byte, 10000),
		errorResp: make(chan []byte, 10000),
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}

			go eppTest.handleClientConnection(conn)
		}
	}()

	return eppTest, nil
}

func (s *EPPTestServer) handleClientConnection(conn net.Conn) {
	defer conn.Close()

	bytesGreeting := []byte(greeting)
	greetBytesLength := uint32(4 + len(bytesGreeting))
	err := binary.Write(conn, binary.BigEndian, greetBytesLength)
	if err != nil {
		fmt.Println("Problem when writing response: " + err.Error())
		return
	}
	if _, err = conn.Write(bytesGreeting); err != nil {
		fmt.Println("Problem when sending response: " + err.Error())
		return
	}

	for {
		var rawResponse int32

		err = binary.Read(conn, binary.BigEndian, &rawResponse)
		if err != nil {
			// Client has closed the connection, so we can close it as well.
			break
		}

		rawResponse -= 4
		if rawResponse < 0 {
			fmt.Println("Problem when reading from client: unexpectedEOF")
			break
		}

		newReq, err := readStreamToBytes(conn, rawResponse)
		if err != nil {
			fmt.Println("Problem when reading client request: " + err.Error())
			break
		}

		expected := <- s.expectedReq
		success := <- s.successResp
		error := <- s.errorResp

		if strings.Contains(string(expected), "REPLACE_REQ_ID") {
			reqId := ""
			splitted := strings.Split(string(newReq), "clTRID")

			if len(splitted) == 3 {
				reg, _ := regexp.Compile("[^A-Z0-9]+")
				reqId = reg.ReplaceAllString(splitted[1], "")
			}

			expected = []byte(strings.Replace(string(expected), "REPLACE_REQ_ID", reqId, 1))
			success = []byte(strings.Replace(string(success), "REPLACE_REQ_ID", reqId, 1))
			error = []byte(strings.Replace(string(error), "REPLACE_REQ_ID", reqId, 1))
		}

		comparison := bytes.Compare(expected, newReq)
		var response []byte
		if comparison == 0 {
			response = success
		} else {
			fmt.Println("Comparison failed, request did not match expected.")
			fmt.Printf("Request:\n%+v\nExpected:\n%+v\n", string(newReq), string(expected))
			response = error
		}

		sendBytesLength := uint32(4 + len(response))
		err = binary.Write(conn, binary.BigEndian, sendBytesLength)
		if err != nil {
			fmt.Println("Problem when writing response: " + err.Error())
			break
		}
		if _, err = conn.Write(response); err != nil {
			fmt.Println("Problem when sending response: " + err.Error())
			break
		}
	}
}

func (s *EPPTestServer) SetupNewResponses(expectedReq, successResp, errorResp string) {
	s.expectedReq <- []byte(expectedReq)
	s.successResp <- []byte(successResp)
	s.errorResp <- []byte(errorResp)
}

func (s *EPPTestServer) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}

	return nil
}

func loadTestServerCerts() (tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair("../../testtmp/testserver.crt", "../../testtmp/testserver.key")
	if err != nil {
		return tls.Certificate{}, nil, errors.Wrap(err, "Have the test certificates been created by running 'make create-test-certs'?")
	}

	caCerts, err := ioutil.ReadFile("../../testtmp/rootCA.crt")
	if err != nil {
		return tls.Certificate{}, nil, errors.Wrap(err, "Is the CA certificate in place?")
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCerts); !ok {
		return tls.Certificate{}, nil, errors.New("Unable to parse given CA certificates.")
	}

	return cert, pool, nil
}

var greeting = `<epp xmlns:obj="urn:ietf:params:xml:ns:obj-1.0"
xmlns="urn:ietf:params:xml:ns:epp-1.0">
<greeting>
 <svID>Ficora EPP Server</svID>
 <svDate>2020-06-20T23:59:59.9720308+02:00</svDate>
 <svcMenu>
 <version>1.0</version>
 <lang>en</lang>
 <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
 <objURI>urn:ietf:params:xml:ns:nsset-1.2</objURI>
 <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
 <objURI>urn:ietf:params:xml:ns:keyset-1.3</objURI>
 <svcExtension>
 <extURI>urn:ietf:params:xml:ns:secDNS-1.1</extURI>
 <extURI>urn:ietf:params:xml:ns:domain-ext-1.0</extURI>
 </svcExtension>
 </svcMenu>
 <dcp>
 <access>
 <personal />
 </access>
 <statement>
 <purpose>
 <prov />
 </purpose>
 <recipient>
 <ours />
 <public />
 </recipient>
 <retention>
 <stated />
 </retention>
 </statement>
 </dcp>
</greeting>
</epp>`

var helloReq = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <hello></hello>
</epp>`

var failedCommand = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2400">
      <msg>Command failed</msg>
    </result>
    <trID>
      <svTRID>sgi4sx2</svTRID>
    </trID>
  </response>
</epp>`