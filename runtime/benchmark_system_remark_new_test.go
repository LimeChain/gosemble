package main

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/benchmarking"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSystemRemark(b *testing.B) {
	size, err := benchmarking.NewLinear(0, uint32(blockLength.Max.Normal))
	assert.NoError(b, err)

	benchmarking.Run(b, "system_remark", func(i *benchmarking.Instance) *benchmarkingtypes.BenchmarkResult {
		// arrange
		message := make([]byte, sc.U32(size.Value()))

		// act
		benchmarkResult, err := i.ExecuteExtrinsic(
			"System.remark",
			sc.NewOption[primitives.RawOrigin](nil),
			&signature.TestKeyringPairAlice,
			message,
		)

		// assert
		assert.NoError(b, err)

		return benchmarkResult
	}, size)
}
