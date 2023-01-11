package utils

import (
	"os"
	"strconv"
)

func GetPort() (int, error) {
	v := os.Getenv("PORT")

	if v == "" {
		v = "9876"
	}

	port, err := strconv.Atoi(v)

	if err != nil {
		return -1, err
	}

	return port, nil
}
