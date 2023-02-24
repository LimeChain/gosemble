package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_UsingEncoded_SignedPayload256(t *testing.T) {
	signedPayload, _ := NewSignedPayload(
		Call{
			CallIndex: CallIndex{
				FunctionIndex: 0,
				ModuleIndex:   0,
			},
			Args: sc.Sequence[sc.U8]{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		SignedExtra{
			Era:   NewEra(ImmortalEra{}),
			Nonce: 0,
			Fee:   0,
		},
	)

	var testExamples = []struct {
		label       string
		input       SignedPayload
		expectation []byte
	}{
		{
			label:       "UsingEncoded SignedPayload()",
			input:       signedPayload,
			expectation: []byte{0x40, 0x87, 0x26, 0xbb, 0xea, 0x99, 0xf8, 0xdf, 0x91, 0x5b, 0xa1, 0x49, 0x59, 0x12, 0x2b, 0x3, 0xd, 0x85, 0x7a, 0xc9, 0x15, 0x73, 0x62, 0x41, 0xde, 0x9c, 0x2d, 0xb0, 0x7c, 0xf8, 0x7, 0x7f},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			enc := sc.SequenceU8ToBytes(testExample.input.UsingEncoded())
			assert.Equal(t, testExample.expectation, enc)
		})
	}
}

func Test_UsingEncoded_SignedPayload257(t *testing.T) {
	signedPayload, _ := NewSignedPayload(
		Call{
			CallIndex: CallIndex{
				FunctionIndex: 0,
				ModuleIndex:   0,
			},
			Args: sc.Sequence[sc.U8]{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		SignedExtra{
			Era:   NewEra(ImmortalEra{}),
			Nonce: 0,
			Fee:   0,
		},
	)

	var testExamples = []struct {
		label       string
		input       SignedPayload
		expectation []byte
	}{
		{
			label:       "UsingEncoded SignedPayload()",
			input:       signedPayload,
			expectation: []byte{0xa4, 0x66, 0x65, 0xf5, 0xb4, 0xd2, 0xe4, 0x5f, 0xb8, 0x27, 0x65, 0x8e, 0xf5, 0xdb, 0x51, 0x5d, 0xb9, 0xcf, 0x28, 0x1a, 0x66, 0x91, 0x3d, 0xc4, 0x13, 0xd8, 0xe4, 0x78, 0xa, 0xc2, 0xcb, 0xdc},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			enc := sc.SequenceU8ToBytes(testExample.input.UsingEncoded())
			assert.Equal(t, testExample.expectation, enc)
		})
	}
}
