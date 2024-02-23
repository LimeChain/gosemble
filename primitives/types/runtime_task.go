package types

import "bytes"

// TODO: implement together with the system.do_task call

type RuntimeTask struct{}

func (t RuntimeTask) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (t RuntimeTask) Bytes() []byte {
	return []byte{}
}

func DecodeRuntimeTask(buffer *bytes.Buffer) (RuntimeTask, error) {
	return RuntimeTask{}, nil
}
