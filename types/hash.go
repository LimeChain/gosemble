package types

import sc "github.com/LimeChain/goscale"

type Hash = sc.FixedSequence[sc.U8]        // size 32
type Blake2bHash = sc.FixedSequence[sc.U8] // size 32
