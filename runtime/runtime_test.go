package main

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/crypto/secp256k1"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
)

const POLKADOT_RUNTIME = "../build/polkadot_runtime-v9400.compact.compressed.wasm"
const NODE_TEMPLATE_RUNTIME = "../build/node_template_runtime.wasm"
const WASM_RUNTIME = "../build/runtime.wasm"

var (
	keySystemHash, _           = common.Twox128Hash([]byte("System"))
	keyAccountHash, _          = common.Twox128Hash([]byte("Account"))
	keyAllExtrinsicsLenHash, _ = common.Twox128Hash([]byte("AllExtrinsicsLen"))
	keyAuraHash, _             = common.Twox128Hash([]byte("Aura"))
	keyAuthoritiesHash, _      = common.Twox128Hash([]byte("Authorities"))
	keyBlockHash, _            = common.Twox128Hash([]byte("BlockHash"))
	keyCurrentSlotHash, _      = common.Twox128Hash([]byte("CurrentSlot"))
	keyDigestHash, _           = common.Twox128Hash([]byte("Digest"))
	keyExecutionPhaseHash, _   = common.Twox128Hash([]byte("ExecutionPhase"))
	keyExtrinsicCountHash, _   = common.Twox128Hash([]byte("ExtrinsicCount"))
	keyExtrinsicIndex          = []byte(":extrinsic_index")
	keyExtrinsicDataHash, _    = common.Twox128Hash([]byte("ExtrinsicData"))
	keyLastRuntime, _          = common.Twox128Hash([]byte("LastRuntimeUpgrade"))
	keyNumberHash, _           = common.Twox128Hash([]byte("Number"))
	keyParentHash, _           = common.Twox128Hash([]byte("ParentHash"))
	keyTimestampHash, _        = common.Twox128Hash([]byte("Timestamp"))
	keyTimestampNowHash, _     = common.Twox128Hash([]byte("Now"))
	keyTimestampDidUpdate, _   = common.Twox128Hash([]byte("DidUpdate"))
	keyBlockWeight, _          = common.Twox128Hash([]byte("BlockWeight"))
)

var (
	parentHash     = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
	stateRoot      = common.MustHexToHash("0xd9e8bf89bda43fb46914321c371add19b81ff92ad6923e8f189b52578074b073")
	extrinsicsRoot = common.MustHexToHash("0x105165e71964828f2b8d1fd89904602cfb9b8930951d87eb249aa2d7c4b51ee7")
	blockNumber    = uint64(1)
	sealDigest     = gossamertypes.SealDigest{
		ConsensusEngineID: gossamertypes.BabeEngineID,
		// bytes for SealDigest that was created in setupHeaderFile function
		Data: []byte{158, 127, 40, 221, 220, 242, 124, 30, 107, 50, 141, 86, 148, 195, 104, 213, 178, 236, 93, 190,
			14, 65, 42, 225, 201, 143, 136, 213, 59, 228, 216, 80, 47, 172, 87, 31, 63, 25, 201, 202, 175, 40, 26,
			103, 51, 25, 36, 30, 12, 80, 149, 166, 131, 173, 52, 49, 98, 4, 8, 138, 54, 164, 189, 134},
	}
)

var (
	invalidTransactionStaleErr             = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	invalidTransactionFutureErr            = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionFuture())
	invalidTransactionBadProofErr          = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
	invalidTransactionExhaustsResourcesErr = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
	unknownTransactionNoUnsignedValidator  = primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	invalidTransactionMandatoryValidation  = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
)

var (
	transactionValidityResultStaleErr, _               = primitives.NewTransactionValidityResult(invalidTransactionStaleErr.(primitives.TransactionValidityError))
	transactionValidityResultFutureErr, _              = primitives.NewTransactionValidityResult(invalidTransactionFutureErr.(primitives.TransactionValidityError))
	transactionValidityResultExhaustsResourcesErr, _   = primitives.NewTransactionValidityResult(invalidTransactionExhaustsResourcesErr.(primitives.TransactionValidityError))
	transactionValidityResultNoUnsignedValidatorErr, _ = primitives.NewTransactionValidityResult(unknownTransactionNoUnsignedValidator.(primitives.TransactionValidityError))
	transactionValidityResultMandatoryValidationErr, _ = primitives.NewTransactionValidityResult(invalidTransactionMandatoryValidation.(primitives.TransactionValidityError))

	dispatchOutcome, _             = primitives.NewDispatchOutcome(nil)
	dispatchOutcomeBadOriginErr, _ = primitives.NewDispatchOutcome(primitives.NewDispatchErrorBadOrigin())

	dispatchOutcomeCustomModuleErr, _ = primitives.NewDispatchOutcome(
		primitives.NewDispatchErrorModule(
			primitives.CustomModuleError{
				Index: BalancesIndex,
				Err:   sc.U32(balances.ErrorInsufficientBalance),
			}))

	dispatchOutcomeExistentialDepositErr, _ = primitives.NewDispatchOutcome(
		primitives.NewDispatchErrorModule(
			primitives.CustomModuleError{
				Index: BalancesIndex,
				Err:   sc.U32(balances.ErrorExistentialDeposit),
			}))

	dispatchOutcomeKeepAliveErr, _ = primitives.NewDispatchOutcome(
		primitives.NewDispatchErrorModule(
			primitives.CustomModuleError{
				Index: BalancesIndex,
				Err:   sc.U32(balances.ErrorKeepAlive),
			}))

	applyExtrinsicResultOutcome, _              = primitives.NewApplyExtrinsicResult(dispatchOutcome)
	applyExtrinsicResultExhaustsResourcesErr, _ = primitives.NewApplyExtrinsicResult(invalidTransactionExhaustsResourcesErr.(primitives.TransactionValidityError))
	applyExtrinsicResultBadOriginErr, _         = primitives.NewApplyExtrinsicResult(dispatchOutcomeBadOriginErr)
	applyExtrinsicResultBadProofErr, _          = primitives.NewApplyExtrinsicResult(invalidTransactionBadProofErr.(primitives.TransactionValidityError))

	applyExtrinsicResultCustomModuleErr, _       = primitives.NewApplyExtrinsicResult(dispatchOutcomeCustomModuleErr)
	applyExtrinsicResultExistentialDepositErr, _ = primitives.NewApplyExtrinsicResult(dispatchOutcomeExistentialDepositErr)
	applyExtrinsicResultKeepAliveErr, _          = primitives.NewApplyExtrinsicResult(dispatchOutcomeKeepAliveErr)
)

func newTestRuntime(t *testing.T) (*wazero_runtime.Instance, *runtime.Storage) {
	runtime := wazero_runtime.NewTestInstanceWithTrie(t, WASM_RUNTIME, trie.NewEmptyTrie())
	return runtime, &runtime.Context.Storage
}

func runtimeMetadata(t *testing.T, instance *wazero_runtime.Instance) *ctypes.Metadata {
	bMetadata, err := instance.Metadata()
	assert.NoError(t, err)

	var decoded []byte
	err = scale.Unmarshal(bMetadata, &decoded)
	assert.NoError(t, err)

	metadata := &ctypes.Metadata{}
	err = codec.Decode(decoded, metadata)
	assert.NoError(t, err)

	return metadata
}

func setStorageAccountInfo(t *testing.T, storage *runtime.Storage, account []byte, freeBalance *big.Int, nonce uint32) (storageKey []byte, info gossamertypes.AccountInfo) {
	accountInfo := gossamertypes.AccountInfo{
		Nonce:       nonce,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(freeBalance),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	aliceHash, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, aliceHash...)
	keyStorageAccount = append(keyStorageAccount, account...)

	bytesStorage, err := scale.Marshal(accountInfo)
	assert.NoError(t, err)

	err = (*storage).Put(keyStorageAccount, bytesStorage)
	assert.NoError(t, err)

	return keyStorageAccount, accountInfo
}

func getQueryInfo(t *testing.T, runtime *wazero_runtime.Instance, extrinsic []byte) primitives.RuntimeDispatchInfo {
	buffer := &bytes.Buffer{}

	buffer.Write(extrinsic)
	err := sc.U32(buffer.Len()).Encode(buffer)
	assert.NoError(t, err)

	bytesRuntimeDispatchInfo, err := runtime.Exec("TransactionPaymentApi_query_info", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesRuntimeDispatchInfo)

	dispatchInfo, err := primitives.DecodeRuntimeDispatchInfo(buffer)
	assert.Nil(t, err)

	return dispatchInfo
}

func timestampExtrinsicBytes(t *testing.T, metadata *ctypes.Metadata, time uint64) []byte {
	call, err := ctypes.NewCall(metadata, "Timestamp.set", ctypes.NewUCompactFromUInt(time))
	assert.NoError(t, err)

	expectedExtrinsic := ctypes.NewExtrinsic(call)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = expectedExtrinsic.Encode(*encoder)
	assert.NoError(t, err)

	return extEnc.Bytes()
}

func signExtrinsicSecp256k1(e *ctypes.Extrinsic, o ctypes.SignatureOptions, keyPair *secp256k1.Keypair) error {
	if e.Type() != ctypes.ExtrinsicVersion4 {
		return fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(), e.Type())
	}

	mb, err := codec.Encode(e.Method)
	if err != nil {
		return err
	}

	era := o.Era
	if !o.Era.IsMortalEra {
		era = ctypes.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := ctypes.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: ctypes.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: o.TransactionVersion,
	}

	b, err := codec.Encode(payload)
	if err != nil {
		return err
	}

	digest := blake2b.Sum256(b)
	signature, err := keyPair.Private().Sign(digest[:])
	if err != nil {
		return err
	}

	signerAddress := blake2b.Sum256(keyPair.Public().Encode())

	signerMultiAddress, err := ctypes.NewMultiAddressFromAccountID(signerAddress[:])
	if err != nil {
		return err
	}

	extSig := ctypes.ExtrinsicSignatureV4{
		Signer:    signerMultiAddress,
		Signature: ctypes.MultiSignature{IsEcdsa: true, AsEcdsa: ctypes.NewEcdsaSignature(signature)},
		Era:       era,
		Nonce:     o.Nonce,
		Tip:       o.Tip,
	}

	e.Signature = extSig

	// mark the extrinsic as signed
	e.Version |= ctypes.ExtrinsicBitSigned

	return nil
}
