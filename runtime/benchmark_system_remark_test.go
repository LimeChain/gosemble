package main

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/benchmarking"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var blockLength, _ = system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)

func BenchmarkSystemRemark(b *testing.B) {
	size, err := benchmarking.NewLinear(0, uint32(blockLength.Max.Normal))
	assert.NoError(b, err)

	benchmarking.Run(b, func(i *benchmarking.Instance) {
		// arrange
		message := make([]byte, sc.U32(size.Value()))

		// act
		err := i.ExecuteExtrinsic(
			"System.remark",
			primitives.NewRawOriginSigned(aliceAccountId),
			message,
		)

		// assert
		assert.NoError(b, err)
	}, size)
}
