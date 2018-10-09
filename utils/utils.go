package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

const ErrorNewUuid = Error("failed to create uuid")

type Error string

func (e Error) Error() string { return string(e) }

// GenerateUUID returns an RFC 4122 compliant UUID or an error.
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	for i := 0; i < 1000; i++ {
		n, err := io.ReadFull(rand.Reader, uuid)
		if n != len(uuid) || err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			uuid[8] = uuid[8]&^0xc0 | 0x80 // variant bits; see section 4.1.1
			uuid[6] = uuid[6]&^0xf0 | 0x40 // version 4 (pseudo-random); see section 4.1.3
			return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
		}
	}
	return "", ErrorNewUuid
}
