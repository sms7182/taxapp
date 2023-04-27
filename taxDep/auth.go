package terminal

import (
	"time"
)

func (t *Terminal) GetToken(taxRawId *uint, taxProcessId *uint, requestUniqueId string, privateKey string) (string, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if len(t.token) > 0 && time.Now().Before(t.exp) {
		return t.token, nil
	}

	packet := t.buildRequestPacket(struct {
		Username string `json:"username"`
	}{
		Username: t.clientID,
	}, "GET_TOKEN", requestUniqueId)

	resp, err := t.transferAPI.SendPacket(taxRawId, taxProcessId, requestUniqueId, packet, "GET_TOKEN", "", false, false, privateKey)
	if err != nil {
		return "", err
	}

	t.token = resp.Result.Data["token"].(string)
	t.exp = time.Now().Add(1 * time.Hour)

	return t.token, nil
}
