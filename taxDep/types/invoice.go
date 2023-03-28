package types

type StandardInvoice struct {
	Header     InvoiceHeader       `json:"header"`
	Body       []InvoiceItem       `json:"body"`
	Payments   []InvoicePayments   `json:"payments"`
	Extensions []InvoiceExtensions `json:"extensions"`
}

type InvoiceHeader struct {
	Indati2m int64   `json:"indati2m,omitempty"`
	Indatim  int64   `json:"indatim,omitempty"`
	Inty     int64   `json:"inty,omitempty"`
	Ft       int64   `json:"ft,omitempty"`
	Inno     int64   `json:"inno,omitempty"`
	Irtaxid  string  `json:"irtaxid,omitempty"`
	Scln     int64   `json:"scln,omitempty"`
	Setm     int64   `json:"setm,omitempty"`
	Tins     string  `json:"tins,omitempty"`
	Cap      float64 `json:"cap,omitempty"`
	Bid      string  `json:"bid,omitempty"`
	Insp     float64 `json:"insp,omitempty"`
	Tvop     float64 `json:"tvop,omitempty"`
	Bpc      string  `json:"bpc,omitempty"`
	Tax17    float64 `json:"tax17,omitempty"`
	Taxid    string  `json:"taxid,omitempty"`
	Inp      int64   `json:"inp,omitempty"`
	Scc      string  `json:"scc,omitempty"`
	Ins      int64   `json:"ins,omitempty"`
	Billid   string  `json:"billid,omitempty"`
	Tprdis   float64 `json:"tprdis,omitempty"`
	Tdis     float64 `json:"tdis"`
	Tadis    float64 `json:"tadis,omitempty"`
	Tvam     float64 `json:"tvam"`
	Todam    float64 `json:"todam"`
	Tbill    float64 `json:"tbill,omitempty"`
	Tob      int64   `json:"tob,omitempty"`
	Tinb     string  `json:"tinb,omitempty"`
	Sbc      string  `json:"sbc,omitempty"`
	Bbc      string  `json:"bbc,omitempty"`
	Bpn      string  `json:"bpn,omitempty"`
	Crn      int64   `json:"crn,omitempty"`
	// Cdcn     string  `json:"cdcn,omitempty"`
	// Cdcd     int64    `json:"cdcd,omitempty"`
	// Tonw     float64 `json:"tonw,omitempty"`
	// Torv     float64 `json:"torv,omitempty"`
	// Tocv     float64 `json:"tocv,omitempty"` // 37
}

type InvoiceItem struct {
	Sstid   string  `json:"sstid,omitempty"`
	Sstt    string  `json:"sstt,omitempty"`
	Mu      int64   `json:"mu,omitempty"`
	Am      float64 `json:"am,omitempty"`
	Fee     float64 `json:"fee,omitempty"`
	Cfee    float64 `json:"cfee,omitempty"`
	Cut     string  `json:"cut,omitempty"`
	Exr     string  `json:"exr,omitempty"`
	Prdis   float64 `json:"prdis,omitempty"`
	Dis     float64 `json:"dis"`
	Adis    float64 `json:"adis,omitempty"`
	Vra     float64 `json:"vra"`
	Vam     float64 `json:"vam"`
	Odt     string  `json:"odt,omitempty"`
	Odr     float64 `json:"odr,omitempty"`
	Odam    float64 `json:"odam,omitempty"`
	Olt     string  `json:"olt,omitempty"`
	Olr     float64 `json:"olr,omitempty"`
	Olam    float64 `json:"olam,omitempty"`
	Consfee float64 `json:"consfee,omitempty"`
	Spro    float64 `json:"spro,omitempty"`
	Bros    float64 `json:"bros,omitempty"`
	Tcpbs   float64 `json:"tcpbs,omitempty"`
	Cop     float64 `json:"cop,omitempty"`
	Bsrn    string  `json:"bsrn,omitempty"`
	Vop     string  `json:"vop,omitempty"`
	Tsstam  float64 `json:"tsstam,omitempty"`
	// Nw      float64 `json:"nw,omitempty"`
	// Ssrv    float64 `json:"ssrv,omitempty"`
	// Sscv    float64 `json:"sscv,omitempty"`
}

type InvoicePayments struct {
	Iinn string  `json:"iinn,omitempty"`
	Acn  string  `json:"acn,omitempty"`
	Trmn string  `json:"trmn,omitempty"`
	Trn  string  `json:"trn,omitempty"`
	Pcn  string  `json:"pcn,omitempty"`
	Pid  string  `json:"pid,omitempty"`
	Pdt  float64 `json:"pdt,omitempty"`
	Pv   float64 `json:"pv,omitempty"`
	Pmt  int64   `json:"pmt,omitempty"`
}

type InvoiceExtensions struct {
}
