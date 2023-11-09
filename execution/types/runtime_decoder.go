package types

import (
	"bytes"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeDecoder interface {
	DecodeBlock(buffer *bytes.Buffer) (primitives.Block, error)
	DecodeUncheckedExtrinsic(buffer *bytes.Buffer) (primitives.UncheckedExtrinsic, error)
	DecodeCall(buffer *bytes.Buffer) (primitives.Call, error)
}

type runtimeDecoder[S primitives.SignerAddress] struct {
	modules []types.Module
	extra   primitives.SignedExtra
}

func NewRuntimeDecoder[S primitives.SignerAddress](modules []types.Module, extra primitives.SignedExtra) RuntimeDecoder {
	return runtimeDecoder[S]{
		modules: modules,
		extra:   extra,
	}
}

func (rd runtimeDecoder[S]) DecodeBlock(buffer *bytes.Buffer) (primitives.Block, error) {
	header, err := primitives.DecodeHeader(buffer)
	if err != nil {
		return nil, err
	}
	compact, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	length := compact.ToBigInt().Int64()
	extrinsics := make([]types.UncheckedExtrinsic, length)

	for i := 0; i < len(extrinsics); i++ {
		extrinsic, err := rd.DecodeUncheckedExtrinsic(buffer)
		if err != nil {
			return nil, err
		}
		extrinsics[i] = extrinsic
	}

	return NewBlock(header, extrinsics), nil
}

func (rd runtimeDecoder[S]) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) (primitives.UncheckedExtrinsic, error) {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLenCompact, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	expectedLength := expectedLenCompact.ToBigInt().Int64()
	beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[primitives.ExtrinsicSignature]
	isSigned := version&ExtrinsicBitSigned != 0
	if isSigned {
		sig, err := primitives.DecodeExtrinsicSignature[S](rd.extra, buffer)
		if err != nil {
			return nil, err
		}
		extSignature = sc.NewOption[primitives.ExtrinsicSignature](sig)
	}

	// Decodes the dispatch call, including its arguments.
	function, err := rd.DecodeCall(buffer)
	if err != nil {
		return nil, err
	}

	afterLength := buffer.Len()

	if int(expectedLength) != beforeLength-afterLength {
		log.Critical("invalid length prefix")
	}

	return NewUncheckedExtrinsic(sc.U8(version), extSignature, function, rd.extra), nil
}

func (rd runtimeDecoder[S]) DecodeCall(buffer *bytes.Buffer) (primitives.Call, error) {
	moduleIndex, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	functionIndex, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	module, err := primitives.GetModule(moduleIndex, rd.modules)
	if err != nil {
		return nil, err
	}

	function, ok := module.Functions()[functionIndex]
	if !ok {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Critical(fmt.Sprintf("function index [%d] for module [%d] not found", functionIndex, moduleIndex))
		log.Critical("function index [" + strconv.Itoa(int(functionIndex)) + "] for module [" + strconv.Itoa(int(moduleIndex)) + "] not found")
	}

	function, err = function.DecodeArgs(buffer)
	if err != nil {
		return nil, err
	}

	return function, nil
}
