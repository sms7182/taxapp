package notify

type NotifyRequest struct {
	From    *string  `json:"from"`
	To      []string `json:"to"`
	Content string   `json:"content"`
	Subject *string  `json:"subject"`

	NotifyType    string   `json:"notifyType"`
	HasAttachment bool     `json:"hasAttachment"`
	Bcc           []string `json:"bcc"`
	Client        *string  `json:"client"`
	ClientName    *string  `json:"clientName"`
	ReplyTo       *string  `json:"replyTo"`
	Requester     *string  `json:"requester"`
}
