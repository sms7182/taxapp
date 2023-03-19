package external

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
