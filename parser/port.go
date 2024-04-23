package parser

import (
	"errors"
	"strconv"
)

func ParsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("invaild port range")
	}
	return port, nil
}
