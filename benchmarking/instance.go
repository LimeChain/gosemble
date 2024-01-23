package benchmarking

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

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
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

var (
	errOnlyOneCall = errors.New("Only one extrinsic or block call is allowed per testFb.")
)

// todo copied from runtime/runtime_test.go
var (
	keySystemHash, _  = common.Twox128Hash([]byte("System"))
	keyAccountHash, _ = common.Twox128Hash([]byte("Account"))
	parentHash        = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
)

type Instance struct {
	// Provides a runtime instance allowing test setup by modifying storage and others
	runtime         *wazero_runtime.Instance
	metadata        *ctypes.Metadata
	version         runtime.Version
	storage         *runtime.Storage
	benchmarkResult *benchmarking.BenchmarkResult
	repeats         int
}

// Creates new benchmarking instance which is used as a param in testFn closure functions
// todo describe repeats param better
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
		return nil, fmt.Errorf("failed to get runtime version: %v", err)
	}

	return &Instance{
		runtime:  runtime,
		metadata: metadata,
		version:  version,
		storage:  &runtime.Context.Storage,
		repeats:  repeats,
	}, nil
}

// Returns Storage instance which can be used to modify the state during benchmark tests
func (i *Instance) Storage() *runtime.Storage {
	return i.storage
}

// Returns runtime instance metadata
func (i *Instance) RuntimeMetadata() *ctypes.Metadata {
	return i.metadata
}

// Returns runtime instance version
func (i *Instance) RuntimeVersion() runtime.Version {
	return i.version
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
func (i *Instance) ExecuteBlock() error {
	return nil
}

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

func accountStorageKey(account []byte) []byte {
	pubKey, _ := common.Blake2b128(account)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, pubKey...)
	keyStorageAccount = append(keyStorageAccount, account...)
	return keyStorageAccount
}
