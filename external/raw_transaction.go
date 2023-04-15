package external

import "tax-management/taxDep/types"

type AfterData struct {
	Username string  `json:"user_name"`
	Tinb     string  `json:"tinb"`
	Tins     string  `json:"tins"`
	Adis     float64 `json:"adis"`
	Am       float64 `json:"am"`
	Dis      float64 `json:"dis"`
	Fee      float64 `json:"fee"`
	Indatim  int64   `json:"indatim"`
	Inp      int64   `json:"inp"`
	Ins      int64   `json:"ins"`
	Inty     int64   `json:"inty"`
	Prdis    float64 `json:"prdis"`
	Setm     int64   `json:"setm"`
	Sstid    string  `json:"sstid"`
	Tadis    float64 `json:"tadis"`
	Tbill    float64 `json:"tbill"`
	Tdis     float64 `json:"tdis"`
	Tob      int64   `json:"tob"`
	Todam    float64 `json:"todam"`
	Tprdis   float64 `json:"tprdis"`
	Trn      string  `json:"trn"`
	Tsstam   float64 `json:"tsstam"`
	Tvam     float64 `json:"tvam"`
	Vam      float64 `json:"vam"`
	Vra      float64 `json:"vra"`
	Cap      float64 `json:"cap"`
	Bid      string  `json:"bid"`
	Insp     float64 `json:"insp"`
	IrTaxid  string  `json:"irtaxid"`
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

type RawTransaction struct {
	After       AfterData   `json:"after"`
	Op          string      `json:"op"`
	Source      SourceData  `json:"source"`
	Transaction interface{} `json:"transaction"`
	TsMs        int64       `json:"ts_ms"`
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
	items := []types.InvoiceItem{
		{
			Sstid:  after.Sstid,
			Am:     after.Am,
			Fee:    after.Fee,
			Prdis:  after.Prdis,
			Dis:    after.Dis,
			Adis:   after.Adis,
			Vra:    after.Vra,
			Vam:    after.Vam,
			Tsstam: after.Tsstam,
		},
	}
	return []types.StandardInvoice{
		{
			Header:     header,
			Body:       items,
			Payments:   nil,
			Extensions: nil,
		},
	}
}
