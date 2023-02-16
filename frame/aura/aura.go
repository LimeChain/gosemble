package aura

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
	"reflect"
)

const MaxAuthorities = 100

var (
	EngineId = []byte{'a', 'u', 'r', 'a'}
)

type Slot = sc.U64

func Authorities() int64 {
	auraHash := hashing.Twox128(constants.KeyAura)
	authoritiesHash := hashing.Twox128(constants.KeyAuthorities)

	authorities := storage.GetBytes(append(auraHash, authoritiesHash...))

	return utils.BytesToOffsetAndSize(authorities)
}

func SlotDuration() int64 {
	return utils.BytesToOffsetAndSize(sc.U64(slotDuration()).Bytes())
}

func OnInitialize() types.Weight {
	slot := currentSlotFromDigests()

	if slot.HasValue {
		newSlot := slot.Value

		auraHash := hashing.Twox128(constants.KeyAura)
		currentSlotHash := hashing.Twox128(constants.KeyCurrentSlot)

		currentSlot := storage.GetDecode(append(auraHash, currentSlotHash...), sc.DecodeU64)

		if currentSlot >= newSlot {
			log.Critical("Slot must increase")
		}

		storage.Set(append(auraHash, currentSlotHash...), newSlot.Bytes())

		totalAuthorities := totalAuthorities()
		if totalAuthorities.HasValue {
			_ = currentSlot % totalAuthorities.Value

			// TODO: implement once  Session module is added
			/*
				if T::DisabledValidators::is_disabled(authority_index as u32) {
							panic!(
								"Validator with index {:?} is disabled and should not be attempting to author blocks.",
								authority_index,
							);
						}
			*/
		}

		// TODO: db weight
		// return T::DbWeight::get().reads_writes(2, 1)
	} else {
		// TODO: db weight
		// return T::DbWeight::get().reads(1)
	}

	return types.Weight{}
}

func OnGenesisSession() {
	// TODO: implement once Session module is added
}

func OnNewSession() {
	// TODO: implement once Session module is added
}

func OnTimestampSet(now sc.U64) {
	slotDuration := slotDuration()
	if slotDuration == 0 {
		log.Critical("Aura slot duration cannot be zero.")
	}

	timestampSlot := now / sc.U64(slotDuration)

	auraHash := hashing.Twox128(constants.KeyAura)
	currentSlotHash := hashing.Twox128(constants.KeyCurrentSlot)
	currentSlot := storage.GetDecode(append(auraHash, currentSlotHash...), sc.DecodeU64)

	if currentSlot != timestampSlot {
		log.Critical("Timestamp slot must match `CurrentSlot`")
	}
}

func currentSlotFromDigests() sc.Option[Slot] {
	systemHash := hashing.Twox128(constants.KeySystem)
	digestHash := hashing.Twox128(constants.KeyDigest)
	digest := storage.GetDecode(append(systemHash, digestHash...), types.DecodeDigest)

	for keyDigest, dig := range digest {
		if keyDigest == types.DigestTypePreRuntime {
			for _, digestItem := range dig {
				if reflect.DeepEqual(sc.FixedSequenceU8ToBytes(digestItem.Engine), EngineId) {
					buffer := &bytes.Buffer{}
					buffer.Write(sc.SequenceU8ToBytes(digestItem.Payload))

					return sc.NewOption[Slot](sc.DecodeU64(buffer))
				}
			}
		}
	}

	return sc.NewOption[Slot](nil)
}

func totalAuthorities() sc.Option[sc.U64] {
	auraHash := hashing.Twox128(constants.KeyAura)
	authoritiesHash := hashing.Twox128(constants.KeyAuthorities)

	// `Compact<u32>` is 5 bytes in maximum.
	data := []byte{0, 0, 0, 0, 0}
	option := storage.Read(append(auraHash, authoritiesHash...), data, 0)

	if !option.HasValue {
		return sc.Option[sc.U64]{}
	}

	length := option.Value
	if length > sc.U32(len(data)) {
		length = sc.U32(len(data))
	}

	buffer := &bytes.Buffer{}
	buffer.Write(data[:length])

	compact := sc.DecodeCompact(buffer)

	return sc.NewOption[sc.U64](sc.U64(compact.ToBigInt().Uint64()))
}

func slotDuration() int {
	return timestamp.MinimumPeriod * 2
}