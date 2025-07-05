package infra

import "crypto/rand"

func GenerateId() string {
	return rand.Text()
}
