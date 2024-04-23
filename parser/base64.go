package parser

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(s string) (string, error) {
	s = strings.TrimSpace(s)
	// url safe
	if strings.Contains(s, "-") || strings.Contains(s, "_") {
		s = strings.Replace(s, "-", "+", -1)
		s = strings.Replace(s, "_", "/", -1)
	}
	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}
	decodeStr, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decodeStr), nil
}
