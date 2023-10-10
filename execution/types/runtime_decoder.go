package types

import (
	"bytes"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeDecoder interface {
	DecodeBlock(buffer *bytes.Buffer) Block
	DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic
	DecodeCall(buffer *bytes.Buffer) primitives.Call
}

type runtimeDecoder struct {
	modules map[sc.U8]Module
	extra   primitives.SignedExtra
}

func NewRuntimeDecoder(modules map[sc.U8]Module, extra primitives.SignedExtra) RuntimeDecoder {
	return runtimeDecoder{modules: modules, extra: extra}
}

func (rd runtimeDecoder) DecodeBlock(buffer *bytes.Buffer) Block {
	header := primitives.DecodeHeader(buffer)

	length := sc.DecodeCompact(buffer).ToBigInt().Int64()
	extrinsics := make([]UncheckedExtrinsic, length)

	for i := 0; i < len(extrinsics); i++ {
		extrinsics[i] = rd.DecodeUncheckedExtrinsic(buffer)
	}

	return Block{
		Header:     header,
		Extrinsics: extrinsics,
	}
}

func (rd runtimeDecoder) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLength := sc.DecodeCompact(buffer).ToBigInt().Int64()
	beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[primitives.ExtrinsicSignature]
	isSigned := version&ExtrinsicBitSigned != 0
	if isSigned {
		extSignature = sc.NewOption[primitives.ExtrinsicSignature](primitives.DecodeExtrinsicSignature(rd.extra, buffer))
	}

	// Decodes the dispatch call, including its arguments.
	function := rd.DecodeCall(buffer)

	afterLength := buffer.Len()

	if int(expectedLength) != beforeLength-afterLength {
		log.Critical("invalid length prefix")
	}

	return NewUncheckedExtrinsic(sc.U8(version), extSignature, function, rd.extra)
}

func (rd runtimeDecoder) DecodeCall(buffer *bytes.Buffer) primitives.Call {
	moduleIndex := sc.DecodeU8(buffer)
	functionIndex := sc.DecodeU8(buffer)

	module, ok := rd.modules[moduleIndex]
	if !ok {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Critical(fmt.Sprintf("module with index [%d] not found", moduleIndex))
		log.Critical("module with index [" + strconv.Itoa(int(moduleIndex)) + "] not found")
	}

	function, ok := module.Functions()[functionIndex]
	if !ok {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Critical(fmt.Sprintf("function index [%d] for module [%d] not found", functionIndex, moduleIndex))
		log.Critical("function index [" + strconv.Itoa(int(functionIndex)) + "] for module [" + strconv.Itoa(int(moduleIndex)) + "] not found")
	}

	function = function.DecodeArgs(buffer)

	return function
}
