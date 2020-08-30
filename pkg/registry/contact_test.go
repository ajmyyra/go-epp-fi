package registry

import (
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"testing"
)

func TestClient_CheckContacts(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedContactCheck, contactCheckResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	checks, err := eppTestClient.CheckContacts("username1", "username2", "username3")
	if err != nil {
		t.Errorf("Contact check failed: %s", err)
	}

	for _, check := range checks {
		if check.Id.Name == "username2" {
			if check.IsAvailable {
				t.Errorf("User %s should be unavailable.", check.Name.Name)
			}
		} else {
			if !check.IsAvailable {
				t.Errorf("User %s should be available: %s", check.Name.Name, check.Reason)
			}
		}
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_GetContact(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedContactInfo, contactInfoResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	info, err := eppTestClient.GetContact("username2")
	if err != nil {
		t.Errorf("Fetching contact failed: %s", err)
	}

	if info.PostalInfo.IsFinnish != 1 {
		t.Errorf("Contact should be Finnish, is the struct parsed correctly?")
	}
	if info.Role != 5 {
		t.Errorf("Contact's role should be registrant (5), not %d", info.Role)
	}

	eppTestServer.SetupNewResponses(expectedContactInfo, contactNotFound, failedCommand)
	if _, err = eppTestClient.GetContact("username2"); err == nil {
		t.Errorf("Fetching nonexistent contact should result in error.")
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_CreateContact(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedContactCreation, contactCreationResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	corporateContact, err := epp.NewBusinessContact(
		5,
		true,
		"Special Test Oy",
		"1881545-1",
		"Testi Test",
		"Vantaa",
		"FI",
		[]string{"Tikkurilantie 1", "3. krs", "Call from the door", "This should not pass."},
		"04230",
		"testi@specialtest.fi",
		"+3585633456",
	)

	if _, err = eppTestClient.CreateContact(corporateContact); err == nil {
		t.Errorf("Too long address format should have caused an error.")
	}

	corporateContact.PostalInfo.Addr.Street = []string{"Tikkurilantie 1", "3. krs"}

	contactId, err := eppTestClient.CreateContact(corporateContact)
	if err != nil {
		t.Errorf("Contact creation failed: %s", err)
	}

	if contactId == "" {
		t.Errorf("No valid contact id received: %s", contactId)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_UpdateContact(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedContactUpdate, successfulCommandResponse, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	corporateContact, err := epp.NewBusinessContact(
		5,
		true,
		"Special Test Oy",
		"1881545-1",
		"Another Testperson",
		"Vantaa",
		"FIN",
		[]string{"Tikkurilantie 1", "5. krs"},
		"04230",
		"another@specialtest.fi",
		"+3585633456",
	)

	if err = eppTestClient.UpdateContact("C575808", corporateContact); err == nil {
		t.Errorf("Erroneous country code should have caused an error.")
	}

	corporateContact.PostalInfo.Addr.Country = "FI"

	if err = eppTestClient.UpdateContact("C575808", corporateContact); err != nil {
		t.Errorf("Contact update failed: %s", err)
	}


	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

func TestClient_DeleteContact(t *testing.T) {
	eppTestServer, eppTestClient, err := initTestServerClient(12004)
	if err != nil {
		t.Fatalf("Error when creating server or client for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestServer.SetupNewResponses(expectedContactDeletion, successfulCommandResponse, contactNotFound)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	if err = eppTestClient.DeleteContact("nope"); err == nil {
		t.Errorf("Deleting nonexisting contact should have caused an error.")
	}

	eppTestServer.SetupNewResponses(expectedContactDeletion, successfulCommandResponse, contactNotFound)

	if err = eppTestClient.DeleteContact("C575808"); err != nil {
		t.Errorf("Contact deletion caused an error: %s", err)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the client connection failed: %s", err)
	}
}

var expectedContactCheck = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <check>
      <contact:check xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
        <contact:id>username1</contact:id>
        <contact:id>username2</contact:id>
        <contact:id>username3</contact:id>
      </contact:check>
    </check>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var contactCheckResponse = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <contact:chkData>
        <contact:cd>
          <contact:id avail="1">username1</contact:id>
        </contact:cd>
        <contact:cd>
          <contact:id avail="0">username2</contact:id>
          <contact:reason>In use</contact:reason>
        </contact:cd>
        <contact:cd>
          <contact:id avail="1">username3</contact:id>
        </contact:cd>
      </contact:chkData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>54322-XYZ</svTRID>
    </trID>
  </response>
</epp>`

var expectedContactInfo = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <info>
      <contact:info xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
        <contact:id>username2</contact:id>
      </contact:info>
    </info>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var contactInfoResponse = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <contact:infData xmlns:contact="urn:ietf:params:xml:ns:contact1.0">
        <contact:id>username2</contact:id>
        <contact:role>5</contact:role>
        <contact:type>1</contact:type>
        <contact:postalInfo type="loc">
          <contact:isFinnish>1</contact:isFinnish>
          <contact:name>HR</contact:name>
          <contact:org>Example Inc.</contact:org>
          <contact:birthDate>2005-04-03</contact:birthDate>
          <contact:identity>123423A123F</contact:identity>
          <contact:registernumber>1234312SFAD-5</contact:registernumber>
          <contact:addr>
            <contact:street>123 Example Dr.</contact:street>
            <contact:street>Suite 100</contact:street>
            <contact:city>Dulles</contact:city>
            <contact:sp>VA</contact:sp>
            <contact:pc>20166-6503</contact:pc>
            <contact:cc>US</contact:cc>
          </contact:addr>
        </contact:postalInfo>
        <contact:voice>+3581231234</contact:voice>
        <contact:email>jdoe@example.com</contact:email>
        <contact:legalemail>jdoe@example.com</contact:legalemail>
        <contact:clID>ClientY</contact:clID>
        <contact:crID>ClientX</contact:crID>
        <contact:crDate>1999-04-03T22:00:00.0Z</contact:crDate>
        <contact:upDate>1999-12-03T09:00:00.0Z</contact:upDate>
        <contact:disclose flag="0">
          <contact:voice/>
          <contact:email/>
          <contact:address/>
        </contact:disclose>
      </contact:infData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>54322-XYZ</svTRID>
    </trID>
  </response>
</epp>`

var expectedContactCreation = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <create>
      <contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
        <contact:role>5</contact:role>
        <contact:type>1</contact:type>
        <contact:postalInfo type="loc">
          <contact:isfinnish>1</contact:isfinnish>
          <contact:name>Testi Test</contact:name>
          <contact:org>Special Test Oy</contact:org>
          <contact:registernumber>1881545-1</contact:registernumber>
          <contact:addr>
            <contact:street>Tikkurilantie 1</contact:street>
            <contact:street>3. krs</contact:street>
            <contact:city>Vantaa</contact:city>
            <contact:pc>04230</contact:pc>
            <contact:cc>FI</contact:cc>
          </contact:addr>
        </contact:postalInfo>
        <contact:voice>+3585633456</contact:voice>
        <contact:email>testi@specialtest.fi</contact:email>
        <contact:legalemail>testi@specialtest.fi</contact:legalemail>
        <contact:disclose flag="0">
          <contact:email>0</contact:email>
          <contact:address>1</contact:address>
        </contact:disclose>
      </contact:create>
    </create>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var contactCreationResponse = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <resData>
      <contact:creData>
        <contact:id>C575808</contact:id>
        <contact:crDate>2020-07-05T23:21:16.0445483+03:00</contact:crDate>
      </contact:creData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>sc3hblz</svTRID>
    </trID>
  </response>
</epp>`

var expectedContactUpdate = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
 <command>
  <update>
   <contact:update xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
    <contact:id>C575808</contact:id>
    <contact:add></contact:add>
    <contact:rem></contact:rem>
    <contact:chg xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
     <contact:role>5</contact:role>
     <contact:type>1</contact:type>
     <contact:postalInfo type="loc">
      <contact:isfinnish>1</contact:isfinnish>
      <contact:name>Another Testperson</contact:name>
      <contact:org>Special Test Oy</contact:org>
      <contact:registernumber>1881545-1</contact:registernumber>
      <contact:addr>
       <contact:street>Tikkurilantie 1</contact:street>
       <contact:street>5. krs</contact:street>
       <contact:city>Vantaa</contact:city>
       <contact:pc>04230</contact:pc>
       <contact:cc>FI</contact:cc>
      </contact:addr>
     </contact:postalInfo>
     <contact:voice>+3585633456</contact:voice>
     <contact:email>another@specialtest.fi</contact:email>
     <contact:legalemail>another@specialtest.fi</contact:legalemail>
     <contact:disclose flag="0">
      <contact:email>0</contact:email>
      <contact:address>1</contact:address>
     </contact:disclose>
    </contact:chg>
   </contact:update>
  </update>
  <clTRID>REPLACE_REQ_ID</clTRID>
 </command>
</epp>`

var expectedContactDeletion = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <delete>
      <contact:delete xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
        <contact:id>C575808</contact:id>
      </contact:delete>
    </delete>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var contactNotFound = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2303">
      <msg>Object does not exist</msg>
    </result>
    <resData>
      <contact:infData>
        <contact:role>0</contact:role>
        <contact:type>0</contact:type>
        <contact:postalInfo>
          <contact:isFinnish>0</contact:isFinnish>
        </contact:postalInfo>
        <contact:voice />
        <contact:disclose>
          <contact:infDataDisclose p10:nil="true" xmlns:p10="http://www.w3.org/2001/XMLSchema-instance" />
        </contact:disclose>
      </contact:infData>
    </resData>
    <trID>
      <clTRID>STXF8</clTRID>
      <svTRID>yckddik</svTRID>
    </trID>
  </response>
</epp>`