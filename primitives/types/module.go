package types

import (
	"errors"
	"strconv"

	sc "github.com/LimeChain/goscale"
)

type Module interface {
	InherentProvider
	DispatchModule
	GetIndex() sc.U8
	Functions() map[sc.U8]Call
	PreDispatch(call Call) (sc.Empty, error)
	ValidateUnsigned(source TransactionSource, call Call) (ValidTransaction, error)
	Metadata(mdGenerator *MetadataTypeGenerator) MetadataModule
}

func GetModule(moduleIndex sc.U8, modules []Module) (Module, error) {
	for _, module := range modules {
		if module.GetIndex() == moduleIndex {
			return module, nil
		}
	}
	return nil, errors.New("module with index [" + strconv.Itoa(int(moduleIndex)) + "] not found.")
}

func MustGetModule(moduleIndex sc.U8, modules []Module) Module {
	for _, module := range modules {
		if module.GetIndex() == moduleIndex {
			return module
		}
	}

	panic("module [" + strconv.Itoa(int(moduleIndex)) + "] not found.")
}
