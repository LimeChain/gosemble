package types

import sc "github.com/LimeChain/goscale"

type UncheckedExtrinsic interface {
	sc.Encodable

	Signature() sc.Option[ExtrinsicSignature]
	Function() Call
	Extra() SignedExtra

	IsSigned() bool
	Check() (CheckedExtrinsic, error)
}
