package registry

import (
	"testing"
)

func TestMsgPollingAndAcknowledging(t *testing.T) {
	eppTestServer, err := createEPPTestServer("127.0.0.1", 12003)
	if err != nil {
		t.Fatalf("Error when creating server for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestClient, err := createEPPTestClient("test", "test123", "127.0.0.1", 12003)
	if err != nil {
		t.Fatalf("Error when creating client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedPollReq, newMessages, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	pollMsg, err := eppTestClient.Poll()
	if err != nil {
		t.Errorf("Polling failed: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedPollAck, successfulAck, failedAck)

	count, err := eppTestClient.PollAck(pollMsg.ID)
	if err != nil {
		t.Errorf("Acknowledging a message failed: %v\n", err)
	}

	if count != 0 {
		t.Errorf("Unexpected number of messages left: %d\n", count)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the connection failed: %v\n", err)
	}
}

func TestEmptyMsgQueueAndWrongMsgId(t *testing.T) {
	eppTestServer, err := createEPPTestServer("127.0.0.1", 12003)
	if err != nil {
		t.Fatalf("Error when creating server for tests: %v\n", err)
	}
	defer eppTestServer.Close()

	eppTestClient, err := createEPPTestClient("test", "test123", "127.0.0.1", 12003)
	if err != nil {
		t.Fatalf("Error when creating client for tests: %v\n", err)
	}

	eppTestServer.SetupNewResponses(expectedPollReq, noNewMessages, failedCommand)

	if err = eppTestClient.Connect(); err != nil {
		t.Fatalf("Connecting failed: %v\n", err)
	}

	pollMsg, err := eppTestClient.Poll()
	if err == nil {
		t.Errorf("Polling should have failed: %v\n", pollMsg)
	}

	eppTestServer.SetupNewResponses(expectedPollAck, successfulAck, failedAck)

	count, err := eppTestClient.PollAck("cabd78dd-a0b0-4fe1-b6d0-abd300229250")
	if err == nil {
		t.Errorf("Acknowledging a message should have failed. Amount of messages left: %d\n", count)
	}

	if count != -1 {
		t.Errorf("Unexpected number of messages left: %d\n", count)
	}

	if err = eppTestClient.Close(); err != nil {
		t.Fatalf("Closing the connection failed: %v\n", err)
	}
}

var expectedPollReq = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <poll op="req"></poll>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var expectedPollAck = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <poll op="ack" msgID="cabd78dd-a0b0-4fe1-b4d0-abd300229250"></poll>
    <clTRID>REPLACE_REQ_ID</clTRID>
  </command>
</epp>`

var newMessages = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1301">
      <msg>Command completed successfully; ack to dequeue</msg>
    </result>
    <msgQ count="1" id="cabd78dd-a0b0-4fe1-b4d0-abd300229250">
      <qDate>2020-06-07T02:05:52.267</qDate>
      <msg>Contact created</msg>
    </msgQ>
    <resData>
      <obj:trnData>
        <obj:name>C574767</obj:name>
      </obj:trnData>
    </resData>
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>rzsg9yi</svTRID>
    </trID>
  </response>
</epp>`

var noNewMessages = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1300">
      <msg>Command completed successfully; no messages</msg>
    </result>
    <msgQ count="0" />
    <resData />
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>86xre4g</svTRID>
    </trID>
  </response>
</epp>`

var successfulAck = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <msgQ count="0" id="cabd78dd-a0b0-4fe1-b4d0-abd300229250" />
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>0ya00cc</svTRID>
    </trID>
  </response>
</epp>`

var failedAck = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:host="urn:ietf:params:xml:ns:host-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xmlns:contact="urn:ietf:params:xml:ns:contact-1.0" xmlns:obj="urn:ietf:params:xml:ns:obj-1.0" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="2303">
      <msg>Object does not exist</msg>
    </result>
    <msgQ count="1" />
    <trID>
      <clTRID>REPLACE_REQ_ID</clTRID>
      <svTRID>0w5twe6</svTRID>
    </trID>
  </response>
</epp>`
