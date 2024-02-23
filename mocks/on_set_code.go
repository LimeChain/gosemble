package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type DefaultOnSetCode struct {
	mock.Mock
}

func (d *DefaultOnSetCode) SetCode(codeBlob sc.Sequence[sc.U8]) error {
	args := d.Called(codeBlob)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(error)
}
