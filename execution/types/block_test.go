package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	expectBlockBytes, _ = hex.DecodeString("3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355c143aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a0004")
)

var (
	mockedExtrinsic = new(mocks.UncheckedExtrinsic)
	extrinsics      = sc.Sequence[primitives.UncheckedExtrinsic]{
		mockedExtrinsic,
	}
)

func Test_Block_New(t *testing.T) {
	expect := block{
		header:     header,
		extrinsics: extrinsics,
	}
	target := setupBlock()

	assert.Equal(t, expect, target)
}

func Test_Block_Encode(t *testing.T) {
	target := setupBlock()
	buffer := &bytes.Buffer{}

	mockedExtrinsic.On("Encode", buffer).Return()

	err := target.Encode(buffer)
	assert.NoError(t, err)

	assert.Equal(t, expectBlockBytes, buffer.Bytes())
}

func Test_Block_Bytes(t *testing.T) {
	target := setupBlock()

	assert.Equal(t, expectBlockBytes, target.Bytes())
}

func Test_Block_Header(t *testing.T) {
	target := setupBlock()

	assert.Equal(t, header, target.Header())
}

func Test_Block_Extrinsics(t *testing.T) {
	target := setupBlock()

	assert.Equal(t, extrinsics, target.Extrinsics())
}

func setupBlock() block {
	target := NewBlock(header, extrinsics)

	return target.(block)
}
