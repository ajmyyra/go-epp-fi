package registry

import "testing"

func TestLoginAndLogout(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12002)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedLogin, successfulLogin, failedLogin)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if err = eppTestClient.Login(); err != nil {
		t.Errorf("Login failed: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedLogout, successfulLogout, failedLogout)

	if err = eppTestClient.Logout(); err != nil {
		t.Errorf("Logout failed: %v\n", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the connection failed: %v\n", err)
	}
}

func TestFailedLogin(t *testing.T) {
	eppTestServer, err := createEPPTestServer("127.0.0.1", 12002)
	if err != nil {
		t.Fatalf("Error when creating server for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestClient, err := createEPPTestClient("test", "wrongPass", "127.0.0.1", 12002)
	if err != nil {
		t.Fatalf("Error when creating client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedLogin, successfulLogin, failedLogin)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if err = eppTestClient.Login(); err == nil {
		t.Error("Login should have failed.")
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the connection failed: %v\n", err)
	}
}

var expectedLogin = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <login>
      <clID>test</clID>
      <pw>test123</pw>
      <options>
        <version>1.0</version>
        <lang>en</lang>
      </options>
      <svcs>
        <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
        <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
        <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
        <svcExtension>
          <extURI>urn:ietf:params:xml:ns:secDNS-1.1</extURI>
          <extURI>urn:ietf:params:xml:ns:domain-ext-1.0</extURI>
        </svcExtension>
      </svcs>
    </login>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var successfulLogin = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>wp3dozy</svTRID>
    </trID>
  </response>
</epp>`

var failedLogin = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2200">
      <msg>Authentication error</msg>
    </result>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>0ovtl0d</svTRID>
    </trID>
  </response>
</epp>`

var expectedLogout = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <logout></logout>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var successfulLogout = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1500">
      <msg>Command completed successfully; ending session</msg>
    </result>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>ltszty7</svTRID>
    </trID>
  </response>
</epp>`

var failedLogout = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="7020">
      <msg>Session idle time exceeded</msg>
    </result>
    <trID />
  </response>
</epp>`
