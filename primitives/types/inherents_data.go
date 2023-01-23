package types

import "bytes"

type InherentsData struct{}

func (i InherentsData) Encode(buffer *bytes.Buffer) {
	panic("not implemented InherentsData Encode")
}

func DecodeInherentsData(buffer *bytes.Buffer) InherentsData {
	panic("not implemented DecodeInherentsData")
}
