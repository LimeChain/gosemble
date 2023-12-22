package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Digest struct {
	sc.Sequence[DigestItem]
}

func NewDigest(items sc.Sequence[DigestItem]) Digest {
	return Digest{items}
}

func DecodeDigest(buffer *bytes.Buffer) (Digest, error) {
	compactSize, err := sc.DecodeCompact[sc.Numeric](buffer)
	if err != nil {
		return Digest{}, err
	}
	size := int(compactSize.ToBigInt().Int64())

	items := sc.Sequence[DigestItem]{}

	for i := 0; i < size; i++ {
		item, err := DecodeDigestItem(buffer)
		if err != nil {
			return Digest{}, err
		}

		items = append(items, item)
	}

	return Digest{items}, nil
}

// PreRuntimes returns a sequence of DigestPreRuntime, containing only DigestItemPreRuntime items
func (d Digest) PreRuntimes() (sc.Sequence[DigestPreRuntime], error) {
	result := sc.Sequence[DigestPreRuntime]{}

	for _, item := range d.Sequence {
		if item.IsPreRuntime() {
			preRuntime, err := item.AsPreRuntime()
			if err != nil {
				return nil, err
			}
			result = append(result, preRuntime)
		}
	}

	return result, nil
}

// OnlyPreRuntimes returns a new Digest, containing only PreRuntime digest items
func (d Digest) OnlyPreRuntimes() Digest {
	items := sc.Sequence[DigestItem]{}
	for _, item := range d.Sequence {
		if item.IsPreRuntime() {
			items = append(items, item)
		}
	}

	return NewDigest(items)
}
