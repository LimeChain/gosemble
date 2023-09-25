package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// RuntimeDbWeight is the weight of database operations that the runtime can invoke.
//
// NOTE: This is currently only measured in computational time, and will probably
// be updated all together once proof size is accounted for.
type RuntimeDbWeight struct {
	Read  sc.U64
	Write sc.U64
}

func (dbw RuntimeDbWeight) Encode(buffer *bytes.Buffer) {
	dbw.Read.Encode(buffer)
	dbw.Write.Encode(buffer)
}

func (dbw RuntimeDbWeight) Bytes() []byte {
	return sc.EncodedBytes(dbw)
}

func (dbw RuntimeDbWeight) Reads(r sc.U64) Weight {
	return WeightFromParts(dbw.Read*r, 0) // saturating_mul
}

func (dbw RuntimeDbWeight) Writes(w sc.U64) Weight {
	return WeightFromParts(dbw.Write*w, 0) // saturating_mul
}

func (dbw RuntimeDbWeight) ReadsWrites(r, w sc.U64) Weight {
	readWeight := dbw.Read * r                        // saturating_mul
	writeWeight := dbw.Write * w                      // saturating_mul
	return WeightFromParts(readWeight+writeWeight, 0) // saturating_add
}
