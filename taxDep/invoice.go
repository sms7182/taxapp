package terminal

import (
	"github.com/google/uuid"
	"tax-management/taxDep/types"
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
