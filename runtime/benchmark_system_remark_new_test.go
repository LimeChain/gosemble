package main

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/benchmarking"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

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
