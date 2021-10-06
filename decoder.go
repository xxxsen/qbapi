package qbapi

import (
	"encoding/json"
	"fmt"
)

type Decoder func([]byte, interface{}) error

var (
	JsonDec = json.Unmarshal
	StrDec  = strDec
)

func strDec(data []byte, v interface{}) error {
	st, ok := v.(*string)
	if !ok {
		return fmt.Errorf("should use string to decode")
	}
	*st = string(data)
	return nil
}
