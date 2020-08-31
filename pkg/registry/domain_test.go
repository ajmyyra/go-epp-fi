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
	defer eppTestServer.Close()

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
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

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
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
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
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainNSUpdate, successfulCommandResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	domainUpdate := epp.NewDomainUpdateNameservers("testdomain2.fi", []string{"ns1.foobar.fi", "ns2.foobar.fi"}, []string{"jill.ns.cloudflare.com", "joe.ns.cloudflare.com"})

	if err := eppTestClient.UpdateDomain(domainUpdate); err != nil {
		t.Errorf("Domain update for nameservers failed: %s", err)
	}

	eppTestServer.SetupNewResponses(expectedDomainTransferKeyUpdate, successfulCommandResponse, failedCommand)

	if domainUpdate, err = epp.NewDomainUpdateSetTransferKey("testdomain2.fi", "invalidKey123"); err == nil {
		t.Errorf("Command should have failed due to invalid key. Received DomainUpdate struct: %+v", domainUpdate)
	}
	domainUpdate, err = epp.NewDomainUpdateSetTransferKey("testdomain2.fi", "fgs+562Fds")
	if err != nil {
		t.Errorf("Received error when asking for transfer key domain update: %s", err)
	}

	if err = eppTestClient.UpdateDomain(domainUpdate); err != nil {
		t.Errorf("Domain update for setting transfer key failed: %s", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_RenewDomain(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainRenewal, successfulDomainRenewal, domainRenewalErrorIncorrectExpiration)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if renewalInfo, err := eppTestClient.RenewDomain("testdomain2.fi", "2020-08-03", 3); err == nil {
		t.Errorf("Domain renewal should have failed with wrong expiration date. Renewal data: %+v", renewalInfo)
	}

	eppTestServer.SetupNewResponses(expectedDomainRenewal, successfulDomainRenewal, domainRenewalErrorIncorrectExpiration)

	renewalInfo, err := eppTestClient.RenewDomain("testdomain2.fi", "2021-08-03", 3)
	if err != nil {
		t.Errorf("Domain renewal failed: %s", err)
	}

	if renewalInfo.Name != "testdomain2.fi" {
		t.Errorf("Wrong domain name: %s", renewalInfo.Name)
	}

	if renewalInfo.ExDate.Before(time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expiration date before expected: %s", renewalInfo.ExDate)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_TransferDomain(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainTransfer, successfulDomainTransfer, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	transfer, err := eppTestClient.TransferDomain("newdomain.fi", "fooBar45+Test", nil)
	if err != nil {
		t.Errorf("Domain transfer failed: %s", err)
	}

	if transfer.Name != "newdomain.fi" {
		t.Errorf("Wrong domain name: %s", transfer.Name)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_UpdateDomainExtensions(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainDNSSecAddition, successfulCommandResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	newDnsSecRecord, err := epp.NewDomainDNSSecRecord(123456, 3, 1, "38EC35D5B3A34B44C39B", 257, 233, 1,">AQPJ////4Q==")

	ext := epp.NewDomainDNSSecUpdateExtension([]epp.DomainDSData{newDnsSecRecord}, nil, true)

	if err = eppTestClient.UpdateDomainExtensions("testdomain2.fi", ext); err != nil {
		t.Errorf("DNSSec update failed: %s", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_DeleteDomain(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12005)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedDomainDeletion, successfulCommandResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if err = eppTestClient.DeleteDomain("testdomain3.fi"); err != nil {
		t.Errorf("Domain deletion failed: %s", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
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
        <domain:crDate>1999-08-03T22:00:00.0Z</domain:crDate>
        <domain:upDate>2020-08-03T09:00:00.0Z</domain:upDate>
        <domain:exDate>2021-08-03T22:00:00.0Z</domain:exDate>
        <domain:trDate>2019-04-08T09:00:00.0Z</domain:trDate>
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

var expectedDomainNSUpdate = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <update>
      <domain:update xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain2.fi</domain:name>
        <domain:add>
          <domain:ns>
            <domain:hostObj>jill.ns.cloudflare.com</domain:hostObj>
            <domain:hostObj>joe.ns.cloudflare.com</domain:hostObj>
          </domain:ns>
        </domain:add>
        <domain:rem>
          <domain:ns>
            <domain:hostObj>ns1.foobar.fi</domain:hostObj>
            <domain:hostObj>ns2.foobar.fi</domain:hostObj>
          </domain:ns>
        </domain:rem>
        <domain:chg></domain:chg>
      </domain:update>
    </update>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var expectedDomainTransferKeyUpdate = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <update>
      <domain:update xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain2.fi</domain:name>
        <domain:add></domain:add>
        <domain:rem></domain:rem>
        <domain:chg>
          <domain:authInfo>
            <domain:pw>fgs+562Fds</domain:pw>
          </domain:authInfo>
        </domain:chg>
      </domain:update>
    </update>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var expectedDomainRenewal = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <renew>
      <domain:renew xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain2.fi</domain:name>
        <domain:curExpDate>2021-08-03</domain:curExpDate>
        <domain:period unit="y">3</domain:period>
      </domain:renew>
    </renew>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var successfulDomainRenewal = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <domain:renData>
        <domain:name>testdomain2.fi</domain:name>
        <domain:exDate>2024-08-03T16:27:27.743</domain:exDate>
      </domain:renData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>cqtn4xv</svTRID>
    </trID>
  </response>
</epp>`

var domainRenewalErrorIncorrectExpiration = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2306">
      <msg>Parameter value policy error: Incorrect expiration date.</msg>
    </result>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>0ze63jl</svTRID>
    </trID>
  </response>
</epp>`

var expectedDomainTransfer = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <transfer op="request">
      <domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>newdomain.fi</domain:name>
        <domain:authInfo>
          <domain:pw>fooBar45+Test</domain:pw>
        </domain:authInfo>
      </domain:transfer>
    </transfer>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var successfulDomainTransfer = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <obj:trnData xmlns:obj="urn:ietf:params:xml:ns:obj">
        <obj:name>newdomain.fi</obj:name>
        <obj:trStatus>Transferred</obj:trStatus>
        <obj:reID>ClientX</obj:reID>
        <obj:reDate>2020-08-03T22:00:00.0Z</obj:reDate>
        <obj:acID>ClientY</obj:acID>
      </obj:trnData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>cqtn4xv</svTRID>
    </trID>
  </response>
</epp>`

var expectedDomainDNSSecAddition = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <update>
      <domain:update xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain2.fi</domain:name>
        <domain:add></domain:add>
        <domain:rem></domain:rem>
        <domain:chg></domain:chg>
      </domain:update>
    </update>
    <extension>
      <secDNS:update xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1">
        <secDNS:rem>
          <secDNS:all>true</secDNS:all>
        </secDNS:rem>
        <secDNS:add>
          <secDNS:dsData>
            <secDNS:keyTag>123456</secDNS:keyTag>
            <secDNS:alg>3</secDNS:alg>
            <secDNS:digestType>1</secDNS:digestType>
            <secDNS:digest>38EC35D5B3A34B44C39B</secDNS:digest>
            <secDNS:keyData>
              <secDNS:flags>257</secDNS:flags>
              <secDNS:protocol>233</secDNS:protocol>
              <secDNS:alg>1</secDNS:alg>
              <secDNS:pubKey>&gt;AQPJ////4Q==</secDNS:pubKey>
            </secDNS:keyData>
          </secDNS:dsData>
        </secDNS:add>
        <secDNS:chg></secDNS:chg>
      </secDNS:update>
    </extension>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var expectedDomainDeletion = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <delete>
      <domain:delete xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>testdomain3.fi</domain:name>
      </domain:delete>
    </delete>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`