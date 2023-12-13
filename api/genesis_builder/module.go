package genesisbuilder

import (
	"bytes"
	"fmt"
	"strings"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GenesisBuilder"
	apiVersion    = 1
)

type GenesisBuilder interface {
	CreateDefaultConfig() ([]byte, error)
	BuildConfig(config []byte) error
}

type Module struct {
	modules  []primitives.Module
	memUtils utils.WasmMemoryTranslator
}

func New(modules []primitives.Module) Module {
	return Module{
		modules:  modules,
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (m Module) CreateDefaultConfig() int64 {
	gcs := []string{}
	for _, m := range m.modules {
		genesisBuilder, ok := m.(GenesisBuilder)
		if !ok {
			continue
		}

		gcBz, err := genesisBuilder.CreateDefaultConfig()
		if err != nil {
			log.Critical(err.Error())
		}

		// gcBz[1:len(gcBz)-1] trims first and last characters which represent start and end of the json
		// CreateDefaultConfig returns a valid json (e.g. {"system":{}}), and here we need it as a json field
		gcs = append(gcs, string(gcBz[1:len(gcBz)-1]))
	}

	gcJson := []byte(fmt.Sprintf("{%s}", strings.Join(gcs, ",")))

	return m.memUtils.BytesToOffsetAndSize(sc.BytesToSequenceU8(gcJson).Bytes())
}

func (m Module) BuildConfig(dataPtr int32, dataLen int32) int64 {
	gcBz := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	gcDecoded, err := sc.DecodeSequence[sc.U8](bytes.NewBuffer(gcBz))
	if err != nil {
		log.Critical(err.Error())
	}

	gcDecodedBytes := sc.SequenceU8ToBytes(gcDecoded)

	for _, m := range m.modules {
		genesisBuilder, ok := m.(GenesisBuilder)
		if !ok {
			continue
		}

		if err := genesisBuilder.BuildConfig(gcDecodedBytes); err != nil {
			log.Critical(err.Error())
		}
	}

	return m.memUtils.BytesToOffsetAndSize([]byte{0})
}

func (m Module) Metadata() primitives.RuntimeApiMetadata {
	// todo metadata
	return primitives.RuntimeApiMetadata{}
}
