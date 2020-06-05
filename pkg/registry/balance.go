package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) Balance() (int, error) {
	balanceReq := epp.APIBalance{}
	balanceReq.Xmlns = epp.EPPNamespace
	balanceReq.Command.ClTRID = createRequestID(reqIDLength)

	balanceData, err := xml.Marshal(balanceReq)
	if err != nil {
		return -1, err
	}

	balanceRawResp, err := s.Send(balanceData)
	if err != nil {
		return -1, err
	}

	var balanceResult epp.APIResult
	if err = xml.Unmarshal(balanceRawResp, &balanceResult); err != nil {
		return -1, err
	}

	if balanceResult.Response.Result.Code != 1000 {
		return -1, errors.New("Wrong return status: " + balanceResult.Response.Result.Msg)
	}

	return balanceResult.Response.ResData.BalanceAmount, nil
}