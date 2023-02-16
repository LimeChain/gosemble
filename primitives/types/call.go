package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Call struct {
	CallIndex CallIndex
	Args      sc.Sequence[sc.U8]
}

// TODO: to have something for testing
// until the actual implementation is done
func NewCall(m string, f string, args sc.Sequence[sc.U8]) Call {
	var c Call

	switch m {
	case "System":
		switch f {
		case "remark":
			c.CallIndex.ModuleIndex = 0
			c.CallIndex.FunctionIndex = 0
			c.Args = args
		}
	default:
		log.Critical("invalid Call type")
	}

	return c
}

func (c Call) Encode(buffer *bytes.Buffer) {
	c.CallIndex.Encode(buffer)
	c.Args.Encode(buffer)
}

func DecodeCall(buffer *bytes.Buffer) Call {
	c := Call{}
	c.CallIndex = DecodeCallIndex(buffer)
	c.Args = sc.DecodeSequence[sc.U8](buffer)
	return c
}

func (c Call) Bytes() []byte {
	return sc.EncodedBytes(c)
}

type CallIndex struct {
	ModuleIndex   sc.U8
	FunctionIndex sc.U8
}

func (ci CallIndex) Encode(buffer *bytes.Buffer) {
	ci.ModuleIndex.Encode(buffer)
	ci.FunctionIndex.Encode(buffer)
}

func DecodeCallIndex(buffer *bytes.Buffer) CallIndex {
	ci := CallIndex{}
	ci.ModuleIndex = sc.DecodeU8(buffer)
	ci.FunctionIndex = sc.DecodeU8(buffer)
	return ci
}

func (ci CallIndex) Bytes() []byte {
	return sc.EncodedBytes(ci)
}

func (c Call) GetDispatchInfo() DispatchInfo {
	// TODO
	// impl<#type_impl_gen> #frame_support::dispatch::GetDispatchInfo for #call_ident<#type_use_gen> #where_clause

	// match *self {
	// 	#(
	// 		Self::#fn_name { #( #args_name_pattern_ref, )* } => {
	// 			// __pallet_base_weight = #fn_weight;

	// 			__pallet_weight = <dyn #frame_support::dispatch::WeighData<( #( & #args_type, )* )>>::weigh_data(&__pallet_base_weight, ( #( #args_name, )* ));

	// 			let __pallet_class = <
	// 				dyn #frame_support::dispatch::ClassifyDispatch<
	// 					( #( & #args_type, )* )
	// 				>
	// 			>::classify_dispatch(&__pallet_base_weight, ( #( #args_name, )* ));

	// 			let __pallet_pays_fee = <
	// 				dyn #frame_support::dispatch::PaysFee<( #( & #args_type, )* )>
	// 			>::pays_fee(&__pallet_base_weight, ( #( #args_name, )* ));

	// 			#frame_support::dispatch::DispatchInfo {
	// 				weight: __pallet_weight,
	// 				class: __pallet_class,
	// 				pays_fee: __pallet_pays_fee,
	// 			}
	// 		},
	// 	)*
	// 	Self::__Ignore(_, _) => unreachable!("__Ignore cannot be used"),
	// }
	return DispatchInfo{}
}

func (c Call) Validate() (ok ValidTransaction, err TransactionValidityError) {
	// TODO
	return ok, err
}

func (c Call) PreDispatch() (ok Pre, err TransactionValidityError) {
	// TODO
	ok = Pre{}
	return ok, err
}

func (c Call) PreDispatchUnsigned() (ok Pre, err TransactionValidityError) {
	// TODO
	ok = Pre{}
	return ok, err
}

func (c Call) Dispatch(i interface{}) (ok PostDispatchInfo, err DispatchError) {
	// TODO
	return ok, err
}
