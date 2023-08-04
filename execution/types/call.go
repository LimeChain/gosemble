package types

import (
	"bytes"
	"strconv"

	sc "github.com/LimeChain/goscale"

	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func DecodeCall(buffer *bytes.Buffer) primitives.Call {
	moduleIndex := sc.DecodeU8(buffer)
	functionIndex := sc.DecodeU8(buffer)

	module, ok := config.Modules[moduleIndex]
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
