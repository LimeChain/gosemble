package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type ExtrinsicInitializer struct {
	mock.Mock
}

func (ex *ExtrinsicInitializer) NewChecked(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) types.CheckedExtrinsic {
	args := ex.Called(signed, function, extra)
	return args.Get(0).(types.CheckedExtrinsic)
}
