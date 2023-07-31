package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ModuleDecoder struct {
	modules map[sc.U8]primitives.Module
}

func NewModuleDecoder(modules map[sc.U8]primitives.Module) ModuleDecoder {
	return ModuleDecoder{modules: modules}
}

func (md ModuleDecoder) DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLength := int(sc.DecodeCompact(buffer).ToBigInt().Int64())
	beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()
	isSigned := version&ExtrinsicBitSigned != 0

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[primitives.ExtrinsicSignature]
	if isSigned {
		extSignature = sc.NewOption[primitives.ExtrinsicSignature](primitives.DecodeExtrinsicSignature(buffer))
	}

	// Decodes the dispatch call, including its arguments.
	function := md.DecodeCall(buffer)

	afterLength := buffer.Len()

	if expectedLength != beforeLength-afterLength {
		log.Critical("invalid length prefix")
	}

	return UncheckedExtrinsic{
		Version:   sc.U8(version),
		Signature: extSignature,
		Function:  function,
	}
}

func (md ModuleDecoder) DecodeCall(buffer *bytes.Buffer) primitives.Call {
	moduleIndex := sc.DecodeU8(buffer)
	functionIndex := sc.DecodeU8(buffer)

	module, ok := md.modules[moduleIndex]
	if !ok {
		log.Critical(fmt.Sprintf("module with index [%d] not found", moduleIndex))
	}

	function, ok := module.Functions()[functionIndex]
	if !ok {
		log.Critical(fmt.Sprintf("function index [%d] for module [%d] not found", functionIndex, moduleIndex))
	}

	function = function.DecodeArgs(buffer)

	return function
}

func (md ModuleDecoder) DecodeBlock(buffer *bytes.Buffer) Block {
	header := primitives.DecodeHeader(buffer)

	size := sc.DecodeCompact(buffer)
	length := size.ToBigInt()
	extrinsics := make([]UncheckedExtrinsic, length.Int64())

	for i := 0; i < len(extrinsics); i++ {
		extrinsics[i] = md.DecodeUncheckedExtrinsic(buffer)
	}

	return Block{
		Header:     header,
		Extrinsics: extrinsics,
	}
}
