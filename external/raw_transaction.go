package external

import "tax-management/taxDep/types"

type AfterData struct {
	Taxid   string `json:"taxid"`
	Tinb    string `json:"tinb"`
	Tins    string `json:"tins"`
	Adis    int64  `json:"adis"`
	Am      int64  `json:"am"`
	Dis     int64  `json:"dis"`
	Fee     int64  `json:"fee"`
	Indatim int64  `json:"indatim"`
	Inp     int64  `json:"inp"`
	Ins     int64  `json:"ins"`
	Inty    int64  `json:"inty"`
	Prdis   int64  `json:"prdis"`
	Setm    int64  `json:"setm"`
	Sstid   string `json:"sstid"`
	Tadis   int64  `json:"tadis"`
	Tbill   int64  `json:"tbill"`
	Tdis    int64  `json:"tdis"`
	Tob     int64  `json:"tob"`
	Todam   int64  `json:"todam"`
	Tprdis  int64  `json:"tprdis"`
	Trn     string `json:"trn"`
	Tsstam  int64  `json:"tsstam"`
	Tvam    int64  `json:"tvam"`
	Vam     int64  `json:"vam"`
	Vra     int64  `json:"vra"`
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

func (r RawTransaction) ToStandardInvoice() []types.StandardInvoice {
	after := r.After
	header := types.InvoiceHeader{
		Indatim: after.Indatim,
		Inty:    after.Inty,
		Setm:    after.Setm,
		Tins:    after.Tins,
		Taxid:   after.Taxid,
		Inp:     after.Inp,
		Ins:     after.Ins,
		Tprdis:  float64(after.Tprdis),
		Tdis:    float64(after.Tdis),
		Tadis:   float64(after.Tadis),
		Tvam:    float64(after.Tvam),
		Todam:   float64(after.Todam),
		Tbill:   float64(after.Tbill),
		Tob:     after.Tob,
		Tinb:    after.Tinb,
	}
	items := []types.InvoiceItem{
		{
			Sstid:  after.Sstid,
			Am:     float64(after.Am),
			Fee:    float64(after.Fee),
			Prdis:  float64(after.Prdis),
			Dis:    float64(after.Dis),
			Adis:   float64(after.Adis),
			Vra:    float64(after.Vra),
			Vam:    float64(after.Vam),
			Tsstam: float64(after.Tsstam),
		},
	}
	return []types.StandardInvoice{
		{
			Header: header,
			Body:   items,
			Payments: []types.InvoicePayments{
				{
					Trn: after.Trn,
				},
			},
			Extensions: nil,
		},
	}
}
