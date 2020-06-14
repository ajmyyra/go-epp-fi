package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) CheckContacts(contacts ...string) ([]epp.ItemCheck, error) {
	reqID := createRequestID(reqIDLength)

	contactCheck := epp.APIContactCheck{}
	contactCheck.Xmlns = epp.EPPNamespace
	contactCheck.Command.Check.ContactCheck.Xmlns = epp.ContactNamespace
	contactCheck.Command.ClTRID = reqID

	contactCheck.Command.Check.ContactCheck.ID = contacts

	checkData, err := xml.MarshalIndent(contactCheck, "", "  ")
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
		if item.ContactId.Avail == 1 {
			item.IsAvailable = true
		}

		checkItems = append(checkItems, item)
	}

	return checkItems, nil
}

func (s *Client) CreateContact(contact epp.ContactInfo) (string, error) {
	reqID := createRequestID(reqIDLength)

	if err := contact.Validate(); err != nil {
		return "", err
	}

	contactCreate := epp.APIContactCreation{}
	contactCreate.Xmlns = epp.EPPNamespace
	contactCreate.Command.ClTRID = reqID

	contact.Xmlns = epp.ContactNamespace
	contactCreate.Command.Create.CreateContact = contact

	createData, err := xml.MarshalIndent(contactCreate, "", "  ")
	if err != nil {
		return "", err
	}

	createRawResp, err := s.Send(createData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return "", err
	}

	var createResult epp.APIResult
	if err = xml.Unmarshal(createRawResp, &createResult); err != nil {
		return "", err
	}

	if createResult.Response.Result.Code != 1000 {
		return "", errors.New("Request failed: " + createResult.Response.Result.Msg)
	}

	contactID := createResult.Response.ResData.CreateData.ID
	s.log.Info("Successfully created a new contact.", "contactID", contactID, "reqId", reqID)

	return contactID, nil
}

func (s *Client) GetContact(contactId string) (epp.ContactResponse, error) {
	reqID := createRequestID(reqIDLength)

	contactInfo := epp.APIContactInfo{}
	contactInfo.Xmlns = epp.EPPNamespace
	contactInfo.Command.Info.ContactInfo.Xmlns = epp.ContactNamespace
	contactInfo.Command.ClTRID = reqID

	contactInfo.Command.Info.ContactInfo.ID = contactId

	checkData, err := xml.MarshalIndent(contactInfo, "", "  ")
	if err != nil {
		return epp.ContactResponse{}, err
	}

	infoRawResp, err := s.Send(checkData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.ContactResponse{}, err
	}

	var infoResp epp.APIContactInfoResponse
	if err = xml.Unmarshal(infoRawResp, &infoResp); err != nil {
		return epp.ContactResponse{}, err
	}

	if infoResp.Response.Result.Code != 1000 {
		return epp.ContactResponse{}, errors.New("Request failed: " + infoResp.Response.Result.Msg)
	}

	return infoResp.Response.ResData.ContactInfo, nil
}

func (s *Client) UpdateContact(contactID string, contact epp.ContactInfo) error {
	reqID := createRequestID(reqIDLength)

	if err := contact.Validate(); err != nil {
		return err
	}

	contactUpdate := epp.APIContactUpdate{}
	contactUpdate.Xmlns = epp.EPPNamespace
	contactUpdate.Command.Update.ContactUpdate.Xmlns = epp.ContactNamespace
	contactUpdate.Command.ClTRID = reqID

	contactUpdate.Command.Update.ContactUpdate.ID = contactID
	contactUpdate.Command.Update.ContactUpdate.Chg = contact

	updateData, err := xml.MarshalIndent(contactUpdate, "", " ")
	if err != nil {
		return err
	}

	updateRawResp, err := s.Send(updateData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var updateResp epp.APIContactInfoResponse
	if err = xml.Unmarshal(updateRawResp, &updateResp); err != nil {
		return err
	}

	if updateResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + updateResp.Response.Result.Msg)
	}

	s.log.Info("Successfully updated contact.", "contactID", contactID, "reqID", reqID)

	return nil
}

func (s *Client) DeleteContact(contactID string) error {
	reqID := createRequestID(reqIDLength)

	contactDelete := epp.APIContactDeletion{}
	contactDelete.Xmlns = epp.EPPNamespace
	contactDelete.Command.Delete.ContactDelete.Xmlns = epp.ContactNamespace
	contactDelete.Command.ClTRID = reqID

	contactDelete.Command.Delete.ContactDelete.ID = contactID

	deleteData, err := xml.MarshalIndent(contactDelete, "", "  ")
	if err != nil {
		return err
	}

	deleteRawResp, err := s.Send(deleteData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var deleteResp epp.APIContactInfoResponse
	if err = xml.Unmarshal(deleteRawResp, &deleteResp); err != nil {
		return err
	}

	if deleteResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + deleteResp.Response.Result.Msg)
	}

	s.log.Info("Successfully deleted contact.", "contactID", contactID, "reqID", reqID)

	return nil
}
