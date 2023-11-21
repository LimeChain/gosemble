package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	postDispatchInfoWithFee = PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](WeightFromParts(1, 2)),
		PaysFee:      PaysYes,
	}

	dispatchError = NewDispatchErrorBadOrigin()

	dispatchErrorWithPostInfo = DispatchErrorWithPostInfo[PostDispatchInfo]{
		PostInfo: postDispatchInfoWithFee,
		Error:    dispatchError,
	}
)

func Test_NewDispatchResult_DispatchError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       sc.Encodable
		expectation DispatchResult
	}{
		{
			label:       "DispatchErrorBadOrigin",
			input:       NewDispatchErrorBadOrigin(),
			expectation: DispatchResult(sc.NewVaryingData(dispatchError)),
		},
		{
			label:       "DispatchErrorWithPostInfo[PostDispatchInfo]",
			input:       dispatchErrorWithPostInfo,
			expectation: DispatchResult(sc.NewVaryingData(dispatchErrorWithPostInfo)),
		},
		{
			label:       "Empty",
			input:       sc.Empty{},
			expectation: DispatchResult(sc.NewVaryingData(sc.Empty{})),
		},
		{
			label:       "nil",
			input:       nil,
			expectation: DispatchResult(sc.NewVaryingData(sc.Empty{})),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			result, err := NewDispatchResult(testExample.input)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_NewDispatchResult_DispatchError_TypeError(t *testing.T) {
	result, err := NewDispatchResult(sc.U8(0))

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchResult' type", err.Error())
	assert.Equal(t, DispatchResult{}, result)
}
