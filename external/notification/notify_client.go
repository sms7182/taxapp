package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"tax-management/notify"
	"time"
)

type NotificationClientImpl struct {
	From      string
	NotifyUrl string
	To        string
}

func (nc NotificationClientImpl) FailedNotify(ctx context.Context, failedBodys []notify.FailedBody, failedTemplate string) (*notify.NotifyResult, error) {
	parsedTemplate, err := template.ParseFiles(failedTemplate)
	if err != nil {
		return nil, err

	}
	rec := strings.Split(nc.To, ",")
	var body bytes.Buffer
	var fmessages []string
	for i := 0; i < len(failedBodys); i++ {
		fm := fmt.Sprintf("Message:%s,Int_Trn:%s", failedBodys[i].FailedMessage, failedBodys[i].Int_Trn)
		fmessages = append(fmessages, fm)
	}
	parsedTemplate.Execute(&body, notify.FailedTaxRequest{
		Faileds: fmessages,
		Message: "Failed Texes",
	})

	requester := "taxManagement"
	subject := "failed taxes"
	nr := notify.NotifyRequest{
		To:         rec,
		Content:    body.String(),
		NotifyType: "email",
		Requester:  &requester,
		Subject:    &subject,
		From:       &nc.From,
	}
	jsonBytes, err := json.Marshal(nr)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jsonBytes)

	request, err := http.NewRequest("POST", nc.NotifyUrl, reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)

	if err != nil || resp.StatusCode == 400 {
		fmt.Printf("Send Post has error %v", err.Error())
		return nil, err

	}
	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Send Post has error %v", err.Error())
		return nil, err
	}
	var resBody notify.NotifyResult
	err = json.Unmarshal(bdy, &resBody)
	if err != nil {
		fmt.Printf("unMarshal has error %v", err.Error())
		return nil, err
	}
	return &resBody, nil
}
