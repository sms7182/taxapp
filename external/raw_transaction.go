package external

import "tax-management/taxDep/types"

type Master struct {
	Username string `json:"user_name"`
	Tinb     string `json:"tinb"`
	Tins     string `json:"tins"`

	Indatim int64 `json:"indatim"`
	Inp     int64 `json:"inp"`
	Ins     int64 `json:"ins"`
	Inty    int64 `json:"inty"`
	Setm    int64 `json:"setm"`

	Tadis  float64 `json:"tadis"`
	Tbill  float64 `json:"tbill"`
	Tdis   float64 `json:"tdis"`
	Tob    int64   `json:"tob"`
	Todam  float64 `json:"todam"`
	Tprdis float64 `json:"tprdis"`
	Trn    string  `json:"trn"`

	Tvam float64 `json:"tvam"`

	Cap     float64  `json:"cap"`
	Bid     string   `json:"bid"`
	Insp    float64  `json:"insp"`
	IrTaxid string   `json:"irtaxid"`
	Detail  []Detail `json:"detail"`
}
type Detail struct {
	Sstid string  `json:"sstid"`
	Am    float64 `json:"am"`
	Fee   float64 `json:"fee"`
	Prdis float64 `json:"prdis"`

	Dis  float64 `json:"dis"`
	Adis float64 `json:"adis"`

	Vra    float64 `json:"vra"`
	Vam    float64 `json:"vam"`
	Tsstam float64 `json:"tsstam"`
}
type SourceData struct {
	Connector string      `json:"connector"`
	DB        string      `json:"db"`
	Lsn       int64       `json:"lsn"`
	Name      string      `json:"name"`
	Schema    string      `json:"schema"`
	Sequence  string      `json:"sequence"`
	Snapshot  string      `json:"snapshot"`
	Table     string      `json:"table"`
	TsMs      int64       `json:"ts_ms"`
	TxID      int64       `json:"txId"`
	Version   string      `json:"version"`
	Xmin      interface{} `json:"xmin"`
}

type CustomerDto struct {
	Token      string `json:"token"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	UserName   string `json:"userName"`
}
type RawTransaction struct {
	After       Master      `json:"after"`
	Op          string      `json:"op"`
	Source      SourceData  `json:"source"`
	Transaction interface{} `json:"transaction"`
	TsMs        int64       `json:"ts_ms"`
	Token       string      `json:"token"`
}

func (r RawTransaction) ToStandardInvoice(taxId string) []types.StandardInvoice {
	after := r.After
	header := types.InvoiceHeader{
		Indatim: after.Indatim,
		Inty:    after.Inty,
		Setm:    after.Setm,
		Tins:    after.Tins,
		Taxid:   taxId,
		Inp:     after.Inp,
		Ins:     after.Ins,
		Tprdis:  after.Tprdis,
		Tdis:    after.Tdis,
		Tadis:   after.Tadis,
		Tvam:    after.Tvam,
		Todam:   after.Todam,
		Tbill:   after.Tbill,
		Tob:     after.Tob,
		Tinb:    after.Tinb,
		Cap:     after.Cap,
		Bid:     after.Bid,
		Insp:    after.Insp,
		Irtaxid: after.IrTaxid,
	}
	var items []types.InvoiceItem
	for i := 0; i < len(after.Detail); i++ {
		detail := after.Detail[i]
		items = append(items, types.InvoiceItem{
			Sstid:  detail.Sstid,
			Am:     detail.Am,
			Fee:    detail.Fee,
			Prdis:  detail.Prdis,
			Dis:    detail.Dis,
			Adis:   detail.Adis,
			Vra:    detail.Vra,
			Vam:    detail.Vam,
			Tsstam: detail.Tsstam,
		})
	}

	return []types.StandardInvoice{
		{
			Header:   header,
			Body:     items,
			Payments: nil,
			//	Extensions: nil,
		},
	}
}
