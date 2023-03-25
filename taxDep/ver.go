package terminal

import (
	"fmt"
	verhoff "github.com/pschlump/verhoeff_algorithm"
	"math"
	"time"
)

func GenerateTaxID(t time.Time, clientID string, serial uint) string {
	daysCount := uint(math.Floor(float64(t.Unix()) / 86400.0))
	return fmt.Sprintf("%s%05X%010X%d", clientID, daysCount, serial, GetVerhoeff(daysCount, clientID, serial))
}

func GetVerhoeff(daysCount uint, clientID string, serial uint) int {
	clientIDNorm := ""
	for _, r := range []byte(clientID) {
		if r >= '0' && r <= '9' {
			clientIDNorm += string(r)
		} else {
			clientIDNorm += fmt.Sprint(r)
		}
	}

	return verhoff.GenerateVerhoeff(fmt.Sprintf("%s%06d%012d", clientIDNorm, daysCount, serial))
}
