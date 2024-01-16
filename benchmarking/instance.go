package benchmarking

import (
	"bytes"
	"fmt"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

// todo copied from runtime/runtime_test.go
var (
	keySystemHash, _  = common.Twox128Hash([]byte("System"))
	keyAccountHash, _ = common.Twox128Hash([]byte("Account"))
	parentHash        = common.MustHexToHash("0x0f6d3477739f8a65886135f58c83ff7c2d4a8300a010dfc8b4c5d65ba37920bb")
)

type Instance struct {
	// Provides a runtime instance allowing test setup by modifying storage and others
	runtime  *wazero_runtime.Instance
	metadata *ctypes.Metadata
	version  runtime.Version
	storage  *runtime.Storage
	repeats  int
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
	accountHash, _ := common.Blake2b128(publicKey)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, accountHash...)
	keyStorageAccount = append(keyStorageAccount, publicKey...)

	bAccountInfo, err := scale.Marshal(accountInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal account info: %v", err)
	}

	if err = (*i.storage).Put(keyStorageAccount, bAccountInfo); err != nil {
		return fmt.Errorf("failed to put account info to storage: %v", err)
	}

	return nil
}

// Executes extrinsic with provided call name.
// Accepts optional param signer, which if provided is used to sign the extrinsic.
// Additionally the method appends the benchmark result to instance.benchmarkResults
func (i *Instance) ExecuteExtrinsic(callName string, origin sc.Option[primitives.RawOrigin], signer *signature.KeyringPair, args ...interface{}) (*benchmarkingtypes.BenchmarkResult, error) {
	extrinsic, err := i.newExtrinsic(callName, signer, args)
	if err != nil {
		return nil, err
	}

	benchmarkConfig := benchmarkingtypes.BenchmarkConfig{
		InternalRepeats: sc.U32(i.repeats),
		Extrinsic:       extrinsic,
		Origin:          origin,
	}

	res, err := i.runtime.Exec("Benchmark_run", benchmarkConfig.Bytes())
	if err != nil {
		return nil, err
	}

	benchmarkResult, err := benchmarkingtypes.DecodeBenchmarkResult(bytes.NewBuffer(res))
	if err != nil {
		return nil, fmt.Errorf("failed to decode benchmark result: %v", err)
	}

	return &benchmarkResult, nil
}

// todo
func (i *Instance) ExecuteBlock() error {
	return nil
}

// Internal method that creates and encodes extrinsic
func (i *Instance) newExtrinsic(callName string, signer *signature.KeyringPair, args []interface{}) (sc.Sequence[sc.U8], error) {
	// Create the call
	call, err := ctypes.NewCall(i.metadata, callName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new call: %v", err)
	}

	// Create the extrinsic
	extrinsic := ctypes.NewExtrinsic(call)

	if signer != nil {
		signatureOptions := ctypes.SignatureOptions{
			BlockHash:          ctypes.Hash(parentHash),
			Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
			GenesisHash:        ctypes.Hash(parentHash),
			Nonce:              ctypes.NewUCompactFromUInt(0),
			SpecVersion:        ctypes.U32(i.version.SpecVersion),
			Tip:                ctypes.NewUCompactFromUInt(0),
			TransactionVersion: ctypes.U32(i.version.TransactionVersion),
		}

		if err = extrinsic.Sign(*signer, signatureOptions); err != nil {
			return nil, fmt.Errorf("failed to sign extrinsic: %v", err)
		}
	}

	// Encode the extrinsic
	encodedExtrinsic := bytes.Buffer{}
	encoder := cscale.NewEncoder(&encodedExtrinsic)
	if err = extrinsic.Encode(*encoder); err != nil {
		return nil, fmt.Errorf("failed to encode extrinsic: %v", err)
	}

	return sc.BytesToSequenceU8(encodedExtrinsic.Bytes()), nil
}
