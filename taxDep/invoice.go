package terminal

import "tax-management/taxDep/types"

func (t *Terminal) SendInvoices(invoices []types.StandardInvoice) (*types.AsyncResponse, error) {
	token, err := t.GetToken()
	if err != nil {
		return nil, err
	}

	packets := make([]types.RequestPacket, len(invoices))

	for i, invoice := range invoices {
		packets[i] = *t.buildRequestPacket(invoice, "INVOICE.V01")
	}

	return t.transferAPI.SendPackets(packets, "normal-enqueue", token, true, true)
}
