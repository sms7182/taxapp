package terminal

import (
	"tax-management/taxDep/types"

	"github.com/google/uuid"
)

func (t *Terminal) SendInvoices(taxRawId *uint, taxProcessId *uint, invoices []types.StandardInvoice, privateKey string, customerid string) (*types.AsyncResponse, error) {
	token, err := t.GetToken(taxRawId, taxProcessId, uuid.NewString(), privateKey)
	if err != nil {
		return nil, err
	}

	packets := make([]types.RequestPacket, len(invoices))

	var requestUniqueId string
	for i, invoice := range invoices {
		requestUniqueId = uuid.NewString()
		packets[i] = *t.buildRequestPacket(invoice, "INVOICE.V01", requestUniqueId)
	}
	reqId := uuid.NewString()
	return t.transferAPI.SendPackets(taxRawId, taxProcessId, reqId, packets, "normal-enqueue", token, true, true, privateKey, customerid)
}

func (t *Terminal) InquiryByReferences(taxRawId *uint, taxProcessId *uint, refs []string, privateKey string) ([]types.InquiryResult, error) {
	token, err := t.GetToken(taxRawId, taxProcessId, uuid.NewString(), privateKey)
	if err != nil {
		return nil, err
	}

	version := "INQUIRY_BY_REFERENCE_NUMBER"

	requestUniqueId := uuid.NewString()
	packet := t.buildRequestPacket(struct {
		Refs []string `json:"referenceNumber"`
	}{
		Refs: refs,
	}, version, requestUniqueId)

	resp, err := t.transferAPI.SendPacketInquiry(taxRawId, taxProcessId, requestUniqueId, packet, version, token, false, false, privateKey)
	if err != nil {
		return nil, err
	}

	return resp.Result.Data, nil
}
