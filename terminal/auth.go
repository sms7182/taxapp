package terminal

import (
	"fmt"
	"tax-management/types"
	"time"
)

func (t *Terminal) GetToken() (string, error) {
	if len(t.token) > 0 && time.Now().Before(t.exp) {
		return t.token, nil
	}

	packet := t.BuildRequestPacket(struct {
		Username string `json:"username"`
	}{
		Username: t.clientID,
	}, "GET_TOKEN")

	resp, err := t.TransferAPI.SendPacket(packet, "GET_TOKEN", nil, false, false)
	if err != nil {
		return "", err
	}
	fmt.Println(resp)
	t.token = (*resp).Result.Data["token"].(string)
	t.exp = time.UnixMilli(int64(resp.Result.Data["expiresIn"].(float64)))

	return t.token, nil
}

func (t *Terminal) ServerInfo() (*types.SyncResponse, error) {
	packet := t.BuildRequestPacket(nil, "GET_SERVER_INFORMATION")
	resp, err := t.TransferAPI.SendPacket(packet, "GET_SERVER_INFORMATION", nil, false, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
