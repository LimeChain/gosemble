package executive

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func InitializeBlock(header types.Header) {
	system.ResetEvents()

	if runtimeUpgrade() {
		// TODO: weight
	}

	system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: weight

	system.NoteFinishedInitialize()
}

func runtimeUpgrade() bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	keyLru := append(systemHash, lastRuntimeUpgradeHash...)
	last := storage.Get(keyLru)

	buf := &bytes.Buffer{}
	buf.Write(last)

	lrupi, err := types.DecodeLastRuntimeUpgradeInfo(buf)
	if err != nil {
		panic(err)
	}

	if constants.RuntimeVersion.SpecVersion > sc.U32(lrupi.SpecVersion.ToBigInt().Int64()) ||
		lrupi.SpecName != constants.RuntimeVersion.SpecName {

		valueLru := append(
			sc.ToCompact(uint64(constants.RuntimeVersion.SpecVersion)).Bytes(),
			constants.RuntimeVersion.SpecName.Bytes()...)
		storage.Set(keyLru, valueLru)

		return true
	}

	return false
}

func extractPreRuntimeDigest(digest types.Digest) types.Digest {
	result := types.Digest{}
	for k, v := range digest {
		if k == types.DigestTypePreRuntime {
			result[k] = v
		}
	}

	return result
}
