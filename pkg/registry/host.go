package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) CheckHosts(hosts ...string) ([]epp.ItemCheck, error) {
	reqID := createRequestID(reqIDLength)

	hostCheck := epp.APIHostCheck{}
	hostCheck.Xmlns = epp.EPPNamespace
	hostCheck.Command.Check.HostCheck.Xmlns = epp.HostNamespace
	hostCheck.Command.ClTRID = reqID

	hostCheck.Command.Check.HostCheck.Name = hosts

	checkData, err := xml.MarshalIndent(hostCheck, "", "  ")
	if err != nil {
		return []epp.ItemCheck{}, err
	}

	checkRawResp, err := s.Send(checkData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return []epp.ItemCheck{}, err
	}

	var checkResult epp.APIResult
	if err = xml.Unmarshal(checkRawResp, &checkResult); err != nil {
		return []epp.ItemCheck{}, err
	}

	if checkResult.Response.Result.Code != 1000 {
		return []epp.ItemCheck{}, errors.New("Request failed: " + checkResult.Response.Result.Msg)
	}

	var checkItems []epp.ItemCheck
	for _, item := range checkResult.Response.ResData.ChkData.Cd {
		if item.Name.Avail == 1 {
			item.IsAvailable = true
		}

		checkItems = append(checkItems, item)
	}

	return checkItems, nil
}

func (s *Client) CreateHost(hostname string, ipAddresses []string) (epp.CreateData, error) {
	reqID := createRequestID(reqIDLength)

	hostCreate := epp.APIHostCreation{}
	hostCreate.Xmlns = epp.EPPNamespace
	hostCreate.Command.Create.HostCreate.Xmlns = epp.HostNamespace
	hostCreate.Command.ClTRID = reqID

	addresses, err := formatHostIPs(ipAddresses)
	if err != nil {
		return epp.CreateData{}, err
	}

	hostCreate.Command.Create.HostCreate.Hostname = hostname
	hostCreate.Command.Create.HostCreate.Addr = addresses

	createData, err := xml.MarshalIndent(hostCreate, "", "  ")
	if err != nil {
		return epp.CreateData{}, err
	}

	createRawResp, err := s.Send(createData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.CreateData{}, err
	}

	var createResp epp.APIResult
	if err = xml.Unmarshal(createRawResp, &createResp); err != nil {
		return epp.CreateData{}, err
	}

	if createResp.Response.Result.Code != 1000 {
		return epp.CreateData{}, errors.New("Request failed: " + createResp.Response.Result.Msg)
	}

	createInfo := createResp.Response.ResData.CreateData

	createInfo.CrDate, err = parseDate(createInfo.RawCrDate)
	if err != nil {
		return epp.CreateData{}, err
	}

	return createInfo, nil
}

func (s *Client) GetHost(host string) (epp.HostInfoResp, error) {
	reqID := createRequestID(reqIDLength)

	hostInfo := epp.APIHostInfo{}
	hostInfo.Xmlns = epp.EPPNamespace
	hostInfo.Command.Info.HostInfo.Xmlns = epp.HostNamespace
	hostInfo.Command.ClTRID = reqID

	hostInfo.Command.Info.HostInfo.Name = host

	infoData, err := xml.MarshalIndent(hostInfo, "", "  ")
	if err != nil {
		return epp.HostInfoResp{}, err
	}

	infoRawResp, err := s.Send(infoData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.HostInfoResp{}, err
	}

	var infoResp epp.APIHostInfoResponse
	if err = xml.Unmarshal(infoRawResp, &infoResp); err != nil {
		return epp.HostInfoResp{}, err
	}

	if infoResp.Response.Result.Code != 1000 {
		return epp.HostInfoResp{}, errors.New("Request failed: " + infoResp.Response.Result.Msg)
	}

	hostnameInfo := infoResp.Response.ResData.HostInfo

	hostnameInfo.CrDate, err = parseDate(hostnameInfo.RawCrDate)
	if err != nil {
		return epp.HostInfoResp{}, err
	}

	hostnameInfo.UpDate, err = parseDate(hostnameInfo.RawUpDate)
	if err != nil {
		return epp.HostInfoResp{}, err
	}

	return hostnameInfo, nil
}

func (s *Client) UpdateHost(hostname string, addIPs, removeIPs []string) error {
	reqID := createRequestID(reqIDLength)

	hostUpdate := epp.APIHostUpdate{}
	hostUpdate.Xmlns = epp.EPPNamespace
	hostUpdate.Command.Update.HostUpdate.Xmlns = epp.HostNamespace
	hostUpdate.Command.ClTRID = reqID

	hostUpdate.Command.Update.HostUpdate.Hostname = hostname

	addedAddresses, err := formatHostIPs(addIPs)
	if err != nil {
		return err
	}
	hostUpdate.Command.Update.HostUpdate.Add.Addr = addedAddresses

	removedAddresses, err := formatHostIPs(removeIPs)
	if err != nil {
		return err
	}
	hostUpdate.Command.Update.HostUpdate.Rem.Addr = removedAddresses

	updateData, err := xml.MarshalIndent(hostUpdate, "", "  ")
	if err != nil {
		return err
	}

	updateRawResp, err := s.Send(updateData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var updateResp epp.APIResult
	if err = xml.Unmarshal(updateRawResp, &updateResp); err != nil {
		return err
	}

	if updateResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + updateResp.Response.Result.Msg)
	}

	return nil
}

func (s *Client) DeleteHost(hostname string) error {
	reqID := createRequestID(reqIDLength)

	hostDelete := epp.APIHostDeletion{}
	hostDelete.Xmlns = epp.EPPNamespace
	hostDelete.Command.Delete.HostDelete.Xmlns = epp.HostNamespace
	hostDelete.Command.ClTRID = reqID

	hostDelete.Command.Delete.HostDelete.Hostname = hostname

	deleteData, err := xml.MarshalIndent(hostDelete, "", "  ")
	if err != nil {
		return err
	}

	deleteRawResp, err := s.Send(deleteData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var deleteResp epp.APIResult
	if err = xml.Unmarshal(deleteRawResp, &deleteResp); err != nil {
		return err
	}

	if deleteResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + deleteResp.Response.Result.Msg)
	}

	return nil
}

func formatHostIPs(rawAddresses []string) ([]epp.HostIPAddress, error) {
	var addresses []epp.HostIPAddress
	for _, rawIp := range rawAddresses {
		ip, err := epp.FormatHostIP(rawIp)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, ip)
	}

	return addresses, nil
}