package terminal

import (
	"tax-management/taxDep/types"

	"github.com/google/uuid"
)

func (t *Terminal) SendInvoices(taxRawId *uint, taxProcessId *uint, invoices []types.StandardInvoice) (*types.AsyncResponse, error) {
	token, err := t.GetToken(taxRawId, taxProcessId, uuid.NewString())
	if err != nil {
		return nil, err
	}

	packets := make([]types.RequestPacket, len(invoices))

	requestUniqueId := uuid.NewString()
	for i, invoice := range invoices {
		packets[i] = *t.buildRequestPacket(invoice, "INVOICE.V01", requestUniqueId)
	}

	return t.transferAPI.SendPackets(taxRawId, taxProcessId, requestUniqueId, packets, "normal-enqueue", token, true, true)
}

func (t *Terminal) InquiryByReferences(refs []string) ([]types.InquiryResult, error) {
	token, err := t.GetToken()
	if err != nil {
		return nil, err
	}

	version := "INQUIRY_BY_REFERENCE_NUMBER"

	packet := t.buildRequestPacket(struct {
		Refs []string `json:"referenceNumber"`
	}{
		Refs: refs,
	}, version)

	resp, err := t.transferAPI.SendPacket(packet, version, token, false, false)
	if err != nil {
		return nil, err
	}

	var inquiryResults []types.InquiryResult
	results := resp.Result.Data.([]any)

	for _, result := range results {
		m := result.(map[string]any)
		ir := types.InquiryResult{
			ReferenceNumber: m["referenceNumber"].(string),
			UID:             m["uid"].(string),
			FiscalID:        m["fiscalId"].(string),
			Status:          m["status"].(string),
			PacketType:      m["packetType"].(string),
		}

		data := m["data"].(map[string]any)

		ir.Data = types.InquiryResultData{
			ConfirmationReferenceID: data["confirmationReferenceId"].(string),
			Error:                   data["error"].([]any),
			Success:                 data["success"].(bool),
		}

		if data["warning"] != nil {
			warning := data["warning"].(map[string]any)
			ir.Data.Warning = types.InquiryDataWarning{
				Code:   warning["code"].(string),
				Detail: warning["detail"].([]any),
				Msg:    warning["msg"].(string),
			}
		}

		inquiryResults = append(inquiryResults, ir)
	}

	return inquiryResults, nil
}
