package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) Poll() (epp.PollMessage, error) {
	reqID := createRequestID(reqIDLength)

	pollReq := epp.APIPoll{}
	pollReq.Xmlns = epp.EPPNamespace
	pollReq.Command.Poll.Op = "req"
	pollReq.Command.ClTRID = reqID

	pollData, err := xml.MarshalIndent(pollReq, "", "  ")
	if err != nil {
		return epp.PollMessage{}, err
	}

	pollRawResp, err := s.Send(pollData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
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
		return epp.PollMessage{}, errors.New("Request failed: " + pollResp.Response.Result.Msg)
	}

	date, err := parseDate(pollResp.Response.MsgQ.RawQDate)
	if err != nil {
		return epp.PollMessage{}, err
	}

	pollResp.Response.MsgQ.QDate = date
	pollResp.Response.MsgQ.Name = pollResp.Response.ResData.TrnData.Name

	return pollResp.Response.MsgQ, nil
}

func (s *Client) PollAck(id string) (int, error) {
	reqID := createRequestID(reqIDLength)

	ackReq := epp.APIPoll{}
	ackReq.Xmlns = epp.EPPNamespace
	ackReq.Command.Poll.Op = "ack"
	ackReq.Command.Poll.MsgID = id
	ackReq.Command.ClTRID = reqID

	ackData, err := xml.MarshalIndent(ackReq, "", "  ")
	if err != nil {
		return -1, err
	}

	ackRawResp, err := s.Send(ackData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return -1, err
	}

	var ackResp epp.APIPollResponse
	if err = xml.Unmarshal(ackRawResp, &ackResp); err != nil {
		return -1, err
	}

	if ackResp.Response.Result.Code != 1000 {
		return -1, errors.New("Request failed: " + ackResp.Response.Result.Msg)
	}

	if ackResp.Response.MsgQ.ID != id {
		return -1, errors.New("Wrong message id acked: " + ackResp.Response.MsgQ.ID)
	}

	messagesLeft := ackResp.Response.MsgQ.Count
	s.log.Debug("Message acknowledged successfully.", "message", id, "messagesLeft", messagesLeft)
	return messagesLeft, nil
}