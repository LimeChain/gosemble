package types

import (
	"bytes"
	"math"
	"math/bits"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

// Era An era to describe the longevity of a transaction.
type Era struct {
	IsImmortal sc.Bool
	EraPeriod  sc.U64
	EraPhase   sc.U64
}

// Period and phase are encoded:
// - The period of validity from the block hash found in the signing material.
// - The phase in the period that this transaction's lifetime begins (and, importantly,
// implies which block hash is included in the signature material). If the `period` is
// greater than 1 << 12, then it will be a factor of the times greater than 1<<12 that
// `period` is.
//
// When used on `FRAME`-based runtimes, `period` cannot exceed `BlockHashCount` parameter
// of `system` module.
//
// E.g. with period == 4:
// 0         10        20        30        40
// 0123456789012345678901234567890123456789012
//              |...|
//    authored -/   \- expiry
// phase = 1
// n = Q(current - phase, period) + phase

// NewMortalEra Create a new era based on a period (which should be a power of two between 4 and 65536
// inclusive) and a block number on which it should start (or, for long periods, be shortly
// after the start).
//
// If using `Era` in the context of `FRAME` runtime, make sure that `period`
// does not exceed `BlockHashCount` parameter passed to `system` module, since that
// prunes old blocks and renders transactions immediately invalid.
func NewMortalEra(period sc.U64, current sc.U64) Era {
	// TODO:
	// period = period.checked_next_power_of_two().unwrap_or(1<<16).clamp(4, 1<<16)
	phase := current % period
	quantizeFactor := sc.Max64(period>>12, 1)
	quantizeFactor = phase / quantizeFactor * quantizeFactor
	return Era{
		IsImmortal: false,
		EraPeriod:  period,
		EraPhase:   quantizeFactor,
	}
}

// The transaction is valid forever. The genesis hash must be present in the signed content.
func NewImmortalEra() Era {
	return Era{IsImmortal: true}
}

func (e Era) Encode(buffer *bytes.Buffer) {
	if e.IsImmortal {
		sc.U8(0).Encode(buffer)
		return
	}

	quantizeFactor := sc.Max64(e.EraPeriod>>12, 1)
	encoded := sc.U16(sc.Clamp(bits.TrailingZeros64(uint64(e.EraPeriod))-1, 1, 15)) | sc.U16((e.EraPhase/quantizeFactor)<<4)
	buffer.Write(encoded.Bytes())
}

func DecodeEra(buffer *bytes.Buffer) (Era, error) {
	firstByte, err := sc.DecodeU8(buffer)
	if err != nil {
		return Era{}, err
	}

	if firstByte == 0 {
		return NewImmortalEra(), nil
	} else {
		nextByte, err := sc.DecodeU8(buffer)
		if err != nil {
			return Era{}, err
		}
		encoded := sc.U64(firstByte) + (sc.U64(nextByte) << 8)
		period := sc.U64(2 << (encoded % (1 << 4)))
		quantizeFactor := sc.Max64(period>>12, 1)
		phase := (encoded >> 4) * quantizeFactor

		if period >= 4 && phase < period {
			return NewMortalEra(period, phase), nil
		} else {
			return Era{}, NewTypeError("Era")
		}
	}
}

func (e Era) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// Get the block number of the start of the era whose properties this object
// describes that `current` belongs to.
func (e Era) Birth(current sc.U64) sc.U64 {
	if e.IsImmortal {
		return sc.U64(0)
	} else {
		period := e.EraPeriod
		phase := e.EraPhase
		return ((((sc.Max64(current, phase)) - phase) / period) * period) + phase
	}
}

// Get the block number of the first block at which the era has ended.
func (e Era) Death(current sc.U64) sc.U64 {
	if e.IsImmortal {
		return sc.U64(math.MaxUint64)
	} else {
		return e.Birth(current) + e.EraPeriod
	}
}

func EraTypeDefinition() sc.Sequence[MetadataDefinitionVariant] {
	result := sc.Sequence[MetadataDefinitionVariant]{
		NewMetadataDefinitionVariant(
			"Immortal",
			sc.Sequence[MetadataTypeDefinitionField]{},
			0,
			"Era.Immortal"),
	}

	// this is necessary since the size of the encoded Mortal variant is `u16`, conditional on
	// the value of the first byte being > 0.
	for i := 1; i <= 255; i++ {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		result = append(result, NewMetadataDefinitionVariant(
			// fmt.Sprintf("Mortal%d", i),
			"Mortal "+strconv.Itoa(int(i)),
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),
			},
			sc.U8(i),
			// fmt.Sprintf("Era.Mortal%d", i),
			"Era.Mortal "+strconv.Itoa(int(i)),
		))
	}

	return result
}
