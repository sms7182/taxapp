package terminal

import (
	"time"
)

func (t *Terminal) GetToken(taxRawId *uint, taxProcessId *uint, requestUniqueId string, privateKey string, userName string) (string, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if len(t.token) > 0 && time.Now().Before(t.exp) {
		return t.token, nil
	}

	packet := t.buildRequestPacket(struct {
		Username string `json:"username"`
	}{
		Username: userName,
	}, "GET_TOKEN", requestUniqueId, userName)

	resp, err := t.transferAPI.SendPacket(taxRawId, taxProcessId, requestUniqueId, packet, "GET_TOKEN", "", false, false, privateKey, userName)
	if err != nil {
		return "", err
	}

	t.token = resp.Result.Data["token"].(string)
	t.exp = time.Now().Add(1 * time.Hour)

	return t.token, nil
}
