package terminal

import (
	"time"
)

func (t *Terminal) GetToken() (string, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if len(t.token) > 0 && time.Now().Before(t.exp) {
		return t.token, nil
	}

	packet := t.buildRequestPacket(struct {
		Username string `json:"username"`
	}{
		Username: t.clientID,
	}, "GET_TOKEN")

	resp, err := t.transferAPI.SendPacket(packet, "GET_TOKEN", "", false, false)
	if err != nil {
		return "", err
	}

	t.token = resp.Result.Data["token"].(string)
	t.exp = time.UnixMilli(int64(resp.Result.Data["expiresIn"].(float64)))

	return t.token, nil
}
