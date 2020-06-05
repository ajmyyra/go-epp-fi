package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

const pollDate = "2006-01-02T15:04:05"

func (s *Client) Poll() (epp.PollMessage, error) {
	pollReq := epp.APIPoll{}
	pollReq.Xmlns = epp.EPPNamespace
	pollReq.Command.Poll.Op = "req"
	pollReq.Command.ClTRID = createRequestID(reqIDLength)

	pollData, err := xml.MarshalIndent(pollReq, "", "  ")
	if err != nil {
		return epp.PollMessage{}, err
	}

	pollRawResp, err := s.Send(pollData)
	if err != nil {
		return epp.PollMessage{}, err
	}

	var pollResp epp.APIPollResponse
	if err = xml.Unmarshal(pollRawResp, &pollResp); err != nil {
		return epp.PollMessage{}, err
	}

	if pollResp.Response.Result.Code == 1300 {
		return epp.PollMessage{}, errors.New("No new messages available.")
	}
	if pollResp.Response.Result.Code != 1301 {
		return epp.PollMessage{}, errors.New(pollResp.Response.Result.Msg)
	}


	date, err := parseDate(pollResp.Response.MsgQ.RawQDate, pollDate)
	if err != nil {
		return epp.PollMessage{}, err
	}

	pollResp.Response.MsgQ.QDate = date
	pollResp.Response.MsgQ.Name = pollResp.Response.ResData.TrnData.Name

	return pollResp.Response.MsgQ, nil
}

func (s *Client) PollAck(id string) error {
	// TODO check id validity before polling

	ackReq := epp.APIPoll{}
	ackReq.Xmlns = epp.EPPNamespace
	ackReq.Command.Poll.Op = "ack"
	ackReq.Command.Poll.MsgID = id
	ackReq.Command.ClTRID = createRequestID(reqIDLength)

	ackData, err := xml.Marshal(ackReq)
	if err != nil {
		return err
	}

	ackRawResp, err := s.Send(ackData)
	if err != nil {
		return err
	}

	var ackResp epp.APIPollResponse
	if err = xml.Unmarshal(ackRawResp, &ackResp); err != nil {
		return err
	}

	if ackResp.Response.Result.Code != 1000 {
		return errors.New("Error: " + ackResp.Response.Result.Msg)
	}

	if ackResp.Response.MsgQ.ID != id {
		return errors.New("Wrong message id acked: " + ackResp.Response.MsgQ.ID)
	}

	// TODO ackResp.Response.MsgQ.Count has number of messages left in queue

	return nil
}