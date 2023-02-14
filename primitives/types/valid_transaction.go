package types

import (
	"math"

	sc "github.com/LimeChain/goscale"
)

// The source of the transaction.
//
// Depending on the source we might apply different validation schemes.
// For instance we can disallow specific kinds of transactions if they were not produced
// by our local node (for instance off-chain workers).
const (
	// Transaction is already included in block.
	//
	// This means that we can't really tell where the transaction is coming from,
	// since it's already in the received block. Note that the custom validation logic
	// using either `Local` or `External` should most likely just allow `InBlock`
	// transactions as well.
	InBlock = sc.U8(iota)

	// Transaction is coming from a local source.
	//
	// This means that the transaction was produced internally by the node
	// (for instance an Off-Chain Worker, or an Off-Chain Call), as opposed
	// to being received over the network.
	Local

	// Transaction has been received externally.
	//
	// This means the transaction has been received from (usually) "untrusted" source,
	// for instance received over the network or RPC.
	External
)

type TransactionSource sc.VaryingData

func NewTransactionSource(value sc.Encodable) TransactionSource {
	// TODO: add validation
	return TransactionSource(sc.NewVaryingData(value))
}

type TransactionPriority = sc.U64
type TransactionLongevity = sc.U64
type TransactionTag = sc.Sequence[sc.U8]

// Contains information concerning a valid transaction.
type ValidTransaction struct {
	// Priority of the transaction.
	//
	// Priority determines the ordering of two transactions that have all
	// their dependencies (required tags) satisfied.
	Priority TransactionPriority

	// Transaction dependencies
	//
	// A non-empty list signifies that some other transactions which provide
	// given tags are required to be included before that one.
	Requires sc.Sequence[TransactionTag]

	// Provided tags
	//
	// A list of tags this transaction provides. Successfully importing the transaction
	// will enable other transactions that depend on (require) those tags to be included as well.
	// Provided and required tags allow Substrate to build a dependency graph of transactions
	// and import them in the right (linear) order.
	Provides sc.Sequence[TransactionTag]

	// Transaction longevity
	//
	// Longevity describes minimum number of blocks the validity is correct.
	// After this period transaction should be removed from the pool or revalidated.
	Longevity TransactionLongevity

	// A flag indicating if the transaction should be propagated to other peers.
	//
	// By setting `false` here the transaction will still be considered for
	// including in blocks that are authored on the current node, but will
	// never be sent to other peers.
	Propagate sc.Bool
}

func DefaultValidTransaction() ValidTransaction {
	return ValidTransaction{
		Priority:  0,
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: TransactionLongevity(math.MaxInt64),
		Propagate: true,
	}
}

// Combine two instances into one, as a best effort. This will take the superset of each of the
// `provides` and `requires` tags, it will sum the priorities, take the minimum longevity and
// the logic *And* of the propagate flags.
func (vt ValidTransaction) CombineWith(other ValidTransaction) ValidTransaction {
	// TODO:
	longevity := sc.U64(math.Min(float64(vt.Longevity), float64(other.Longevity)))

	return ValidTransaction{
		// Priority: vt.Priority.saturating_add(other.Priority),
		// Requires: {
		// 	vt.Requires.append(other.Requires),
		// 	vt.Requires,
		// },
		// Provides: {
		// 	vt.Provides.append(other.Provides),
		// 	vt.Provides,
		// },
		Longevity: longevity,
		Propagate: vt.Propagate && other.Propagate,
	}
}
