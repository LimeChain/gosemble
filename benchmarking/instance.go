package benchmarking

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

var (
	errOnlyOneCall = errors.New("Only one extrinsic or block call is allowed per testFb.")
)

var (
	blockNumber = uint(1)
	dateTime    = uint64(time.Date(2023, time.January, 2, 3, 4, 5, 6, time.UTC).UnixMilli())
	parentHash  = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
)

var (
	keySystemHash, _  = common.Twox128Hash([]byte("System"))
	keyAccountHash, _ = common.Twox128Hash([]byte("Account"))
)

type Instance struct {
	// Provides a runtime instance allowing test setup by modifying storage and others
	runtime         *wazero_runtime.Instance
	metadata        *ctypes.Metadata
	storage         *runtime.Storage
	version         runtime.Version
	benchmarkResult *benchmarking.BenchmarkResult
	repeats         int
}

// Creates new benchmarking instance which is used as a param in testFn closure functions
func newBenchmarkingInstance(runtime *wazero_runtime.Instance, repeats int) (*Instance, error) {
	bMetadata, err := runtime.Metadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get runtime metadata: %v", err)
	}

	var metadataDecoded []byte
	if err = scale.Unmarshal(bMetadata, &metadataDecoded); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	metadata := &ctypes.Metadata{}
	if err = codec.Decode(metadataDecoded, metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %v", err)
	}

	version, err := runtime.Version()
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	return &Instance{
		runtime:  runtime,
		metadata: metadata,
		storage:  &runtime.Context.Storage,
		version:  version,
		repeats:  repeats,
	}, nil
}

// Returns Storage instance which can be used to modify the state during benchmark tests
func (i *Instance) Storage() *runtime.Storage {
	return i.storage
}

// Sets the specified account info for the specified public key
func (i *Instance) SetAccountInfo(publicKey []byte, accountInfo gossamertypes.AccountInfo) error {
	bAccountInfo, err := scale.Marshal(accountInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal account info: %v", err)
	}

	if err = (*i.storage).Put(accountStorageKey(publicKey), bAccountInfo); err != nil {
		return fmt.Errorf("failed to put account info to storage: %v", err)
	}

	return nil
}

func (i *Instance) GetAccountInfo(publicKey []byte) (gossamertypes.AccountInfo, error) {
	bytesStorage := (*i.storage).Get(accountStorageKey(publicKey))

	accountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(0)),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	err := scale.Unmarshal(bytesStorage, &accountInfo)

	return accountInfo, err
}

// Executes extrinsic with provided call name.
// Accepts optional param signer, which if provided is used to sign the extrinsic.
// Additionally the method appends the benchmark result to instance.benchmarkResults
func (i *Instance) ExecuteExtrinsic(callName string, origin primitives.RawOrigin, args ...interface{}) error {
	if i.benchmarkResult != nil {
		return errOnlyOneCall
	}

	extrinsic, err := i.newExtrinsic(callName, args)
	if err != nil {
		return err
	}

	benchmarkConfig := benchmarkingtypes.BenchmarkConfig{
		InternalRepeats: sc.U32(i.repeats),
		Extrinsic:       extrinsic,
		Origin:          origin,
	}

	res, err := i.runtime.Exec("Benchmark_run", benchmarkConfig.Bytes())
	if err != nil {
		return err
	}

	benchmarkResult, err := benchmarkingtypes.DecodeBenchmarkResult(bytes.NewBuffer(res))
	if err != nil {
		return fmt.Errorf("failed to decode benchmark result: %v", err)
	}

	i.benchmarkResult = &benchmarkResult

	return nil
}

// todo
// func (i *Instance) ExecuteBlock() error {
// 	return nil
// }

// Internal method that creates and encodes extrinsic
func (i *Instance) newExtrinsic(callName string, args []interface{}) (sc.Sequence[sc.U8], error) {
	// Create the call
	call, err := ctypes.NewCall(i.metadata, callName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new call: %v", err)
	}

	// Create the extrinsic
	extrinsic := ctypes.NewExtrinsic(call)

	// Encode the extrinsic
	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	if err = extrinsic.Encode(*encoder); err != nil {
		return nil, fmt.Errorf("failed to encode extrinsic: %v", err)
	}

	return sc.BytesToSequenceU8(encodedExtrinsic.Bytes()), nil
}

func (i *Instance) newSignedExtrinsic(signer signature.KeyringPair, signatureOptions ctypes.SignatureOptions, callName string, args ...interface{}) ([]byte, error) {
	// Create the call
	call, err := ctypes.NewCall(i.metadata, callName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new call: %v", err)
	}

	// Create the extrinsic
	extrinsic := ctypes.NewExtrinsic(call)

	// Sign the extrinsic
	err = extrinsic.Sign(signer, signatureOptions)
	if err != nil {
		return nil, err
	}

	// Encode the extrinsic
	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	if err := extrinsic.Encode(*encoder); err != nil {
		return nil, err
	}

	return encodedExtrinsic.Bytes(), nil
}

func (i *Instance) BuildGenesisConfig() error {
	genesisConfig := []byte("{\"system\":{},\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]},\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]},\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1000000000000000000]]},\"transactionPayment\":{\"multiplier\":\"2\"}}")

	result, err := i.runtime.Exec("GenesisBuilder_build_config", sc.BytesToSequenceU8(genesisConfig).Bytes())
	if err != nil {
		return err
	}

	if !reflect.DeepEqual([]byte{0}, result) {
		return fmt.Errorf("failed to build genesis config: [%v]", result)
	}

	return nil
}

func accountStorageKey(account []byte) []byte {
	pubKey, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, pubKey...)
	keyStorageAccount = append(keyStorageAccount, account...)
	return keyStorageAccount
}

func timestampInherentData(dateTime uint64) ([]byte, error) {
	idata := gossamertypes.NewInherentData()
	err := idata.SetInherent(gossamertypes.Timstap0, dateTime)
	if err != nil {
		return nil, err
	}

	ienc, err := idata.Encode()
	if err != nil {
		return nil, err
	}

	return ienc, err
}
