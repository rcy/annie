package model

import (
	"fmt"
	"goirc/internal/idstr"
	"os"
)

func (n Note) Link() (string, error) {
	str, err := idstr.Encode(n.ID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", os.Getenv("ROOT_URL"), str), nil
}
