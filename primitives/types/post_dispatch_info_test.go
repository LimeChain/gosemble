package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesPostDispatchInfo, _ = hex.DecodeString("01040800")
)

var (
	dispatchInfo = &DispatchInfo{
		Weight:  WeightFromParts(3, 4),
		Class:   NewDispatchClassMandatory(),
		PaysFee: PaysNo,
	}
	postDispatchInfo = PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](WeightFromParts(1, 2)),
		PaysFee:      PaysYes,
	}
)

func Test_PostDispatchInfo_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := postDispatchInfo.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesPostDispatchInfo, buffer.Bytes())
}

func Test_DecodePostDispatchInfo(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesPostDispatchInfo)

	result, err := DecodePostDispatchInfo(buffer)
	assert.NoError(t, err)

	assert.Equal(t, postDispatchInfo, result)
}

func Test_PostDispatchInfo_Bytes(t *testing.T) {
	result := postDispatchInfo.Bytes()

	assert.Equal(t, expectBytesPostDispatchInfo, result)
}

func Test_PostDispatchInfo_CalcUnspent_NoWeight(t *testing.T) {
	target := PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](nil),
		PaysFee:      PaysNo,
	}
	result := target.CalcUnspent(dispatchInfo)

	assert.Equal(t, WeightFromParts(0, 0), result)
}

func Test_PostDispatchInfo_CalcUnspent(t *testing.T) {
	result := postDispatchInfo.CalcUnspent(dispatchInfo)

	assert.Equal(t, WeightFromParts(2, 2), result)
}

func Test_PostDispatchInfo_CalcActualWeight(t *testing.T) {
	result := postDispatchInfo.CalcActualWeight(dispatchInfo)

	assert.Equal(t, postDispatchInfo.ActualWeight.Value, result)
}

func Test_PostDispatchInfo_CalcActualWeight_NoWeight(t *testing.T) {
	target := PostDispatchInfo{
		ActualWeight: sc.NewOption[Weight](nil),
		PaysFee:      PaysNo,
	}
	result := target.CalcActualWeight(dispatchInfo)

	assert.Equal(t, dispatchInfo.Weight, result)
}

func Test_PostDispatchInfo_Pays_Yes(t *testing.T) {
	dispatchInfo := &DispatchInfo{
		PaysFee: PaysYes,
	}
	result := postDispatchInfo.Pays(dispatchInfo)

	assert.Equal(t, PaysYes, result)
}

func Test_PostDispatchInfo_Pays_No(t *testing.T) {
	result := postDispatchInfo.Pays(dispatchInfo)

	assert.Equal(t, PaysNo, result)
}
