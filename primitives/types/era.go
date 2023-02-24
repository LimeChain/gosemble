package types

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// The transaction is valid forever. The genesis hash must be present in the signed content.
type ImmortalEra = sc.Empty

// Period and phase are encoded:
// - The period of validity from the block hash found in the signing material.
// - The phase in the period that this transaction's lifetime begins (and, importantly,
// implies which block hash is included in the signature material). If the `period` is
// greater than 1 << 12, then it will be a factor of the times greater than 1<<12 that
// `period` is.
//
// When used on `FRAME`-based runtimes, `period` cannot exceed `BlockHashCount` parameter
// of `system` module.
type MortalEra struct {
	EraPeriod sc.U64
	EraPhase  sc.U64
}

func (e MortalEra) Encode(buffer *bytes.Buffer) {
	e.EraPeriod.Encode(buffer)
	e.EraPhase.Encode(buffer)
}

func DecodeMortalEra(buffer *bytes.Buffer) MortalEra {
	return MortalEra{
		EraPeriod: sc.DecodeU64(buffer),
		EraPhase:  sc.DecodeU64(buffer),
	}
}

func (e MortalEra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// An era to describe the longevity of a transaction.
type Era sc.VaryingData

// E.g. with period == 4:
// 0         10        20        30        40
// 0123456789012345678901234567890123456789012
//              |...|
//    authored -/   \- expiry
// phase = 1
// n = Q(current - phase, period) + phase

// Create a new era based on a period (which should be a power of two between 4 and 65536
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
	quantizeFactor := (period >> 12).Max(1)
	quantizeFactor = phase / quantizeFactor * quantizeFactor
	return Era(sc.NewVaryingData(MortalEra{EraPeriod: period, EraPhase: quantizeFactor}))
}

func NewEra(value sc.Encodable) Era {
	switch value.(type) {
	case ImmortalEra, MortalEra:
		return Era(sc.NewVaryingData(value))
	default:
		log.Critical("invalid Era type")
	}

	panic("unreachable")
}

func (e Era) Encode(buffer *bytes.Buffer) {
	switch e[0].(type) {
	case ImmortalEra:
		sc.U8(0).Encode(buffer)
	case MortalEra:
		sc.U8(1).Encode(buffer)
		e[0].Encode(buffer)
	default:
		log.Critical("invalid Era type")
	}
}

func DecodeEra(buffer *bytes.Buffer) Era {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return Era(sc.NewVaryingData(ImmortalEra{}))
	case 1:
		return Era(sc.NewVaryingData(DecodeMortalEra(buffer)))
	default:
		log.Critical("invalid Era type")
	}

	panic("unreachable")
}

func (e Era) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e Era) IsImmortalEra() sc.Bool {
	switch e[0].(type) {
	case ImmortalEra:
		return true
	default:
		return false
	}
}

func (e Era) AsImmortalEra() ImmortalEra {
	if e.IsImmortalEra() {
		return e[0].(ImmortalEra)
	} else {
		log.Critical("not a ImmortalEra type")
	}

	panic("unreachable")
}

func (e Era) IsMortalEra() sc.Bool {
	switch e[0].(type) {
	case MortalEra:
		return true
	default:
		return false
	}
}

func (e Era) AsMortalEra() MortalEra {
	if e.IsMortalEra() {
		return e[0].(MortalEra)
	} else {
		log.Critical("not a MortalEra type")
	}

	panic("unreachable")
}

// Get the block number of the start of the era whose properties this object
// describes that `current` belongs to.
func (e Era) Birth(current sc.U64) sc.U64 {
	if e.IsImmortalEra() {
		return sc.U64(0)
	}

	if e.IsMortalEra() {
		period := e.AsMortalEra().EraPeriod
		phase := e.AsMortalEra().EraPhase
		return (current.Max(phase)-phase)/period*period + phase
	}

	log.Critical("invalid era")
	panic("unreachable")
}

// Get the block number of the first block at which the era has ended.
func (e Era) Death(current sc.U64) sc.U64 {
	if e.IsImmortalEra() {
		return sc.U64(math.MaxUint64)
	}

	if e.IsMortalEra() {
		return e.Birth(current) + e.AsMortalEra().EraPeriod
	}

	log.Critical("invalid era")
	panic("unreachable")
}
