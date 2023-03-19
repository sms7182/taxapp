package types

type StandardInvoice struct {
	Header     InvoiceHeader       `json:"header"`
	Body       []InvoiceItem       `json:"body"`
	Payments   []InvoicePayments   `json:"payments"`
	Extensions []InvoiceExtensions `json:"extensions"`
}

type InvoiceHeader struct {
	Indati2m uint    `json:"indati2m"`
	Indatim  uint    `json:"indatim"`
	Inty     uint    `json:"inty"`
	Ft       uint    `json:"ft"`
	Inno     uint    `json:"inno"`
	Irtaxid  string  `json:"irtaxid"`
	Scln     uint    `json:"scln"`
	Setm     uint    `json:"setm"`
	Tins     string  `json:"tins"`
	Cap      float64 `json:"cap"`
	Bid      string  `json:"bid"`
	Insp     float64 `json:"insp"`
	Tvop     float64 `json:"tvop"`
	Bpc      string  `json:"bpc"`
	Tax17    float64 `json:"tax17"`
	Taxid    string  `json:"taxid"`
	Inp      uint    `json:"inp"`
	Scc      string  `json:"scc"`
	Ins      uint    `json:"ins"`
	Billid   string  `json:"billid"`
	Tprdis   float64 `json:"tprdis"`
	Tdis     float64 `json:"tdis"`
	Tadis    float64 `json:"tadis"`
	Tvam     float64 `json:"tvam"`
	Todam    float64 `json:"todam"`
	Tbill    float64 `json:"tbill"`
	Tob      uint    `json:"tob"`
	Tinb     string  `json:"tinb"`
	Sbc      string  `json:"sbc"`
	Bbc      string  `json:"bbc"`
	Bpn      string  `json:"bpn"`
	Crn      uint    `json:"crn"`
	// Cdcn     string  `json:"cdcn"`
	// Cdcd     uint    `json:"cdcd"`
	// Tonw     float64 `json:"tonw"`
	// Torv     float64 `json:"torv"`
	// Tocv     float64 `json:"tocv"` // 37
}

type InvoiceItem struct {
	Sstid   string  `json:"sstid"`
	Sstt    string  `json:"sstt"`
	Mu      uint    `json:"mu"`
	Am      float64 `json:"am"`
	Fee     float64 `json:"fee"`
	Cfee    float64 `json:"cfee"`
	Cut     string  `json:"cut"`
	Exr     string  `json:"exr"`
	Prdis   float64 `json:"prdis"`
	Dis     float64 `json:"dis"`
	Adis    float64 `json:"adis"`
	Vra     float64 `json:"vra"`
	Vam     float64 `json:"vam"`
	Odt     string  `json:"odt"`
	Odr     float64 `json:"odr"`
	Odam    float64 `json:"odam"`
	Olt     string  `json:"olt"`
	Olr     float64 `json:"olr"`
	Olam    float64 `json:"olam"`
	Consfee float64 `json:"consfee"`
	Spro    float64 `json:"spro"`
	Bros    float64 `json:"bros"`
	Tcpbs   float64 `json:"tcpbs"`
	Cop     float64 `json:"cop"`
	Bsrn    string  `json:"bsrn"`
	Vop     string  `json:"vop"`
	Tsstam  float64 `json:"tsstam"`
	// Nw      float64 `json:"nw"`
	// Ssrv    float64 `json:"ssrv"`
	// Sscv    float64 `json:"sscv"`
}

type InvoicePayments struct {
	Iinn string  `json:"iinn"`
	Acn  string  `json:"acn"`
	Trmn string  `json:"trmn"`
	Trn  string  `json:"trn"`
	Pcn  string  `json:"pcn"`
	Pid  string  `json:"pid"`
	Pdt  float64 `json:"pdt"`
	Pv   float64 `json:"pv"`
	Pmt  uint    `json:"pmt"`
}

type InvoiceExtensions struct {
}
