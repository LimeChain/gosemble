package types

import (
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Module interface {
	InherentProvider
	DispatchModule
	GetIndex() sc.U8
	Functions() map[sc.U8]Call
	PreDispatch(call Call) (sc.Empty, TransactionValidityError)
	ValidateUnsigned(source TransactionSource, call Call) (ValidTransaction, TransactionValidityError)
	Metadata() (sc.Sequence[MetadataType], MetadataModule)
}

func GetModule(moduleIndex sc.U8, modules []Module) (Module, bool) {
	for _, module := range modules {
		if module.GetIndex() == moduleIndex {
			return module, true
		}
	}

	return nil, false
}

func MustGetModule(moduleIndex sc.U8, modules []Module) Module {
	for _, module := range modules {
		if module.GetIndex() == moduleIndex {
			return module
		}
	}

	log.Critical("module [" + strconv.Itoa(int(moduleIndex)) + "] not found.")

	panic("unreachable")
}
