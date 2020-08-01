package registry

import (
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"testing"
	"time"
)

func TestClient_CheckDomains(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedDomainCheck, domainCheckResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	checks, err := eppTestClient.CheckDomains("testdomain1.fi", "testdomain2.fi", "testdomain3.fi")
	if err != nil {
		t.Errorf("Domain check failed: %s", err)
	}

	for _, check := range checks {
		if check.Name.Name == "testdomain2.fi" {
			if check.IsAvailable {
				t.Errorf("Domain %s should be unavailable.", check.Name.Name)
			}
		} else {
			if !check.IsAvailable {
				t.Errorf("Domain %s should be available: %s", check.Name.Name, check.Reason)
			}
		}
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_GetDomain(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12006)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedDomainInfo, domainInfoResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	info, err := eppTestClient.GetDomain("testdomain2.fi")
	if err != nil {
		t.Errorf("Fetching domain failed: %s", err)
	}

	if info.Name != "testdomain2.fi" {
		t.Errorf("Wrong domain name: %s", info.Name)
	}
	if len(info.Ns.HostObj) != 2 {
		t.Errorf("Wrong amount of name servers: %d", len(info.Ns.HostObj))
	}
	if len(info.DsData) != 2 {
		t.Errorf("Wrong amount of DNSsec objects: %d", len(info.DsData))
	}
	if info.DsData[0].Alg != 3 || info.DsData[0].KeyData.PubKey != "AQPJ////4Q==" {
		t.Errorf("Malformed DNSsec object: %+v", info.DsData[0])
	}

	eppTestServer.SetupNewResponses(expectedDomainInfo, domainNotFound, failedCommand)
	if nonexistent, err := eppTestClient.GetDomain("testdomain2.fi"); err == nil {
		t.Errorf("Fetching nonexistent domain should result in error: %+v", nonexistent)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_CreateDomain(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainCreation, domainCreationResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	domainDetails := epp.NewDomainDetails("testdomain3.co.uk", 2, "TST1234", []string{"ns1.testhosting.fi", "ns2.testhosting.fi"})

	if err = domainDetails.Validate(); err == nil {
		t.Errorf("Domainobject validation should have returned an error.")
	}

	if _, err = eppTestClient.CreateDomain(domainDetails); err == nil {
		t.Errorf("Too long address format should have caused an error.")
	}

	domainDetails.Name = "testdomain3.fi"

	details, err := eppTestClient.CreateDomain(domainDetails)
	if err != nil {
		t.Errorf("Domain creation failed: %s", err)
	}

	if details.Name != "testdomain3.fi" {
		t.Errorf("Wrong domain name: %s", details.Name)
	}
	if !details.ExDate.After(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Not created for two years, expires at: %s", details.ExDate)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_UpdateDomain(t *testing.T) {

}

func TestClient_RenewDomain(t *testing.T) {

}

func TestClient_TransferDomain(t *testing.T) {

}

func TestClient_UpdateDomainExtensions(t *testing.T) {

}

func TestClient_DeleteDomain(t *testing.T) {

}

var expectedDomainCheck = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <check>
      <domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain1.fi</domain:name>
        <domain:name>testdomain2.fi</domain:name>
        <domain:name>testdomain3.fi</domain:name>
      </domain:check>
    </check>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var domainCheckResponse = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:cd>
          <domain:name avail="1">testdomain1.fi</domain:name>
        </domain:cd>
        <domain:cd>
          <domain:name avail="0">testdomain2.fi</domain:name>
          <domain:reason>In use</domain:reason>
        </domain:cd>
        <domain:cd>
          <domain:name avail="1">testdomain3.fi</domain:name>
        </domain:cd>
      </domain:chkData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>54322-XYZ</svTRID>
    </trID>
  </response>
</epp>`

var expectedDomainInfo = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <info>
      <domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name hosts="all">testdomain2.fi</domain:name>
      </domain:info>
    </info>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var domainInfoResponse = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain2.fi</domain:name>
        <domain:registrylock>1</domain:registrylock>
        <domain:autorenew>1</domain:autorenew>
        <domain:autorenewDate>2018-09-25T12:11:29.433</domain:autorenewDate>
        <domain:status s="Granted"/>
        <domain:registrant>TST1234</domain:registrant>
        <domain:contact type="admin">C2000</domain:contact>
        <domain:contact type="tech">C4000</domain:contact>
        <domain:ns>
          <domain:hostObj>ns1.example.com</domain:hostObj>
          <domain:hostObj>ns1.example.net</domain:hostObj>
        </domain:ns>
        <domain:clID>ClientX</domain:clID>
        <domain:crID>ClientY</domain:crID>
        <domain:crDate>1999-04-03T22:00:00.0Z</domain:crDate>
        <domain:upDate>1999-12-03T09:00:00.0Z</domain:upDate>
        <domain:exDate>2005-04-03T22:00:00.0Z</domain:exDate>
        <domain:trDate>2000-04-08T09:00:00.0Z</domain:trDate>
        <domain:authInfo>
          <domain:pw>2fooBAR</domain:pw>
        </domain:authInfo>
        <domain:dsData>
          <domain:keyTag>12345</domain:keyTag>
          <domain:alg>3</domain:alg>
          <domain:digestType>1</domain:digestType>
          <domain:digest>38EC35D5B3A34B33C99B</domain:digest>
          <domain:keyData>
            <domain:flags>257</domain:flags>
            <domain:protocol>233</domain:protocol>
            <domain:alg>1</domain:alg>
            <domain:pubKey>AQPJ////4Q==</domain:pubKey>
          </domain:keyData>
        </domain:dsData>
        <domain:dsData>
          <domain:keyTag>12345</domain:keyTag>
          <domain:alg>3</domain:alg>
          <domain:digestType>1</domain:digestType>
          <domain:digest>38EC35D5B3A34B33C99B</domain:digest>
          <domain:keyData>
            <domain:flags>257</domain:flags>
            <domain:protocol>233</domain:protocol>
            <domain:alg>1</domain:alg>
            <domain:pubKey>AQPJ////4Q==</domain:pubKey>
          </domain:keyData>
        </domain:dsData>
      </domain:infData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>54322-XYZ</svTRID>
    </trID>
  </response>
</epp>`

var expectedDomainCreation = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <command>
    <create>
      <domain:create xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain3.fi</domain:name>
        <domain:period unit="y">2</domain:period>
        <domain:ns>
          <domain:hostObj>ns1.testhosting.fi</domain:hostObj>
          <domain:hostObj>ns2.testhosting.fi</domain:hostObj>
        </domain:ns>
        <domain:registrant>TST1234</domain:registrant>
      </domain:create>
    </create>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var domainCreationResponse = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <domain:creData>
        <domain:name>testdomain3.fi</domain:name>
        <domain:crDate>2020-08-01T16:27:27.743</domain:crDate>
        <domain:exDate>2022-08-01T16:27:27.743</domain:exDate>
      </domain:creData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>5e4pwd2</svTRID>
    </trID>
  </response>
</epp>`

var domainCreationFailure = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2302">
      <msg>Object exists</msg>
    </result>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>ha8ajqk</svTRID>
    </trID>
  </response>
</epp>`

var domainNotFound = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2303">
      <msg>Object does not exist</msg>
    </result>
    <resData>
      <domain:infData>
        <domain:registrylock>0</domain:registrylock>
        <domain:autorenew>0</domain:autorenew>
        <domain:autorenewDate>0001-01-01T00:00:00</domain:autorenewDate>
        <domain:status />
        <domain:crDate>0001-01-01T00:00:00</domain:crDate>
        <domain:upDate>0001-01-01T00:00:00</domain:upDate>
        <domain:exDate>0001-01-01T00:00:00</domain:exDate>
        <domain:trDate>0001-01-01T00:00:00</domain:trDate>
        <domain:authInfo />
      </domain:infData>
    </resData>
    <extension>
      <infData xmlns="urn:ietf:params:xml:ns:secDNS-1.1" />
      <deletiondate xmlns="urn:ietf:params:xml:ns:domain-ext-1.0">
        <schedule>
          <delDate>0001-01-01T00:00:00</delDate>
        </schedule>
      </deletiondate>
    </extension>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>cqtn4xv</svTRID>
    </trID>
  </response>
</epp>`