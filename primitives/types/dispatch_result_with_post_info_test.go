package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	dispatchInfoWithFee = DispatchInfo{
		Weight:  WeightFromParts(6, 7),
		Class:   NewDispatchClassNormal(),
		PaysFee: PaysYes,
	}

	postDispatchInfoOk = PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](WeightFromParts(2, 3)),
		PaysFee:      PaysYes,
	}

	postDispatchInfoErr = PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](WeightFromParts(4, 5)),
		PaysFee:      PaysNo,
	}
)

func Test_ExtractActualWeight_DispatchResultOk(t *testing.T) {
	dispatchResultWithPostInfo := DispatchResultWithPostInfo[PostDispatchInfo]{
		HasError: false,
		Ok:       postDispatchInfoOk,
		Err: DispatchErrorWithPostInfo[PostDispatchInfo]{
			PostInfo: postDispatchInfoErr,
			Error:    NewDispatchErrorBadOrigin(),
		},
	}

	result := ExtractActualWeight(&dispatchResultWithPostInfo, &dispatchInfoWithFee)

	assert.Equal(t, WeightFromParts(2, 3), result)
}

func Test_ExtractActualWeight_DispatchResultErr(t *testing.T) {
	dispatchResultWithPostInfo := DispatchResultWithPostInfo[PostDispatchInfo]{
		HasError: true,
		Ok:       postDispatchInfoOk,
		Err: DispatchErrorWithPostInfo[PostDispatchInfo]{
			PostInfo: postDispatchInfoErr,
			Error:    NewDispatchErrorBadOrigin(),
		},
	}

	result := ExtractActualWeight(&dispatchResultWithPostInfo, &dispatchInfoWithFee)

	assert.Equal(t, WeightFromParts(4, 5), result)
}

func Test_ExtractActualPaysFee_DispatchResultOk(t *testing.T) {
	dispatchResultWithPostInfo := DispatchResultWithPostInfo[PostDispatchInfo]{
		HasError: false,
		Ok:       postDispatchInfoOk,
		Err: DispatchErrorWithPostInfo[PostDispatchInfo]{
			PostInfo: postDispatchInfoErr,
			Error:    NewDispatchErrorBadOrigin(),
		},
	}

	result := ExtractActualPaysFee(&dispatchResultWithPostInfo, &dispatchInfoWithFee)

	assert.Equal(t, PaysYes, result)
}

func Test_ExtractActualPaysFee_DispatchResultErr(t *testing.T) {
	dispatchResultWithPostInfo := DispatchResultWithPostInfo[PostDispatchInfo]{
		HasError: true,
		Ok:       postDispatchInfoOk,
		Err: DispatchErrorWithPostInfo[PostDispatchInfo]{
			PostInfo: postDispatchInfoErr,
			Error:    NewDispatchErrorBadOrigin(),
		},
	}

	result := ExtractActualPaysFee(&dispatchResultWithPostInfo, &dispatchInfoWithFee)

	assert.Equal(t, PaysNo, result)
}
