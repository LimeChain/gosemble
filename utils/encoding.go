package utils

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

func EncodeEach(buffer *bytes.Buffer, encodables ...sc.Encodable) error {
	for _, encodable := range encodables {
		err := encodable.Encode(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
