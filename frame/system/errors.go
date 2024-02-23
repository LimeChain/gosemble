package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// System module errors.
const (
	ErrorInvalidSpecName sc.U8 = iota
	ErrorSpecVersionNeedsToIncrease
	ErrorFailedToExtractRuntimeVersion
	ErrorNonDefaultComposite
	ErrorNonZeroRefCount
	ErrorCallFiltered
	ErrorInvalidTask
	ErrorFailedTask
	ErrorNothingAuthorized
	ErrorUnauthorized
)

func NewDispatchErrorInvalidSpecName(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorInvalidSpecName),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorSpecVersionNeedsToIncrease(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorSpecVersionNeedsToIncrease),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorFailedToExtractRuntimeVersion(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorFailedToExtractRuntimeVersion),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorNonDefaultComposite(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorNonDefaultComposite),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorNonZeroRefCount(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorNonZeroRefCount),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorCallFiltered(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorCallFiltered),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorInvalidTask(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorInvalidTask),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorFailedTask(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorFailedTask),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorNothingAuthorized(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorNothingAuthorized),
		Message: sc.NewOption[sc.Str](nil),
	})
}

func NewDispatchErrorUnauthorized(moduleId sc.U8) primitives.DispatchError {
	return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorUnauthorized),
		Message: sc.NewOption[sc.Str](nil),
	})
}
