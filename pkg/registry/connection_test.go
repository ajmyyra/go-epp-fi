package registry

import (
	"testing"
)

func TestConnectionAndHello(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12001)
	if err != nil {
		t.Fatalf("Error when creating server for tests: %v\n", err)
	}
	defer eppTestServer.Close()

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