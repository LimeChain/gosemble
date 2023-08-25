package types

import (
	"bytes"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ModuleDecoder[N sc.Numeric] struct {
	modules map[sc.U8]Module[N]
	extra   primitives.SignedExtra
}

func NewModuleDecoder[N sc.Numeric](modules map[sc.U8]Module[N], extra primitives.SignedExtra) ModuleDecoder[N] {
	return ModuleDecoder[N]{modules: modules, extra: extra}
}

func (md ModuleDecoder[N]) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLength := int(sc.To[sc.U64](sc.U128(sc.DecodeCompact(buffer))))
	beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()
	isSigned := version&ExtrinsicBitSigned != 0

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[primitives.ExtrinsicSignature]
	if isSigned {
		extSignature = sc.NewOption[primitives.ExtrinsicSignature](primitives.DecodeExtrinsicSignature(md.extra, buffer))
	}

	// Decodes the dispatch call, including its arguments.
	function := md.DecodeCall(buffer)

	afterLength := buffer.Len()

	if expectedLength != beforeLength-afterLength {
		log.Critical("invalid length prefix")
	}

	return NewUncheckedExtrinsic(sc.U8(version), extSignature, function, md.extra)
}

func (md ModuleDecoder[N]) DecodeCall(buffer *bytes.Buffer) primitives.Call {
	moduleIndex := sc.DecodeU8(buffer)
	functionIndex := sc.DecodeU8(buffer)

	module, ok := md.modules[moduleIndex]
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

func (md ModuleDecoder[N]) DecodeBlock(buffer *bytes.Buffer) Block[N] {
	header := primitives.DecodeHeader[N](buffer)

	length := sc.To[sc.U64](sc.U128(sc.DecodeCompact(buffer)))
	extrinsics := make([]UncheckedExtrinsic, length)

	for i := 0; i < len(extrinsics); i++ {
		extrinsics[i] = md.DecodeUncheckedExtrinsic(buffer)
	}

	return Block[N]{
		Header:     header,
		Extrinsics: extrinsics,
	}
}
