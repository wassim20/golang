package constants

//create a regex color to check with

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (u *UUIDArray) Scan(src interface{}) error {
	var str string

	switch src := src.(type) {
	case []byte:
		str = string(src)
	case string:
		str = src
	default:
		return errors.New("Incompatible type for UUIDArray")
	}

	// Trim the curly braces
	str = str[1 : len(str)-1]

	uuidStrs := strings.Split(str, ",")
	for _, s := range uuidStrs {
		uuid, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		*u = append(*u, uuid)
	}

	return nil
}

func (u UUIDArray) Value() (driver.Value, error) {
	return u, nil
}
