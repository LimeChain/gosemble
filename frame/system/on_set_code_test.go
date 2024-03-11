package system

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	onSetCode defaultOnSetCode
)

var (
	systemModule *mocks.SystemModule
)

func setupDefaultOnSetCode() {
	systemModule = new(mocks.SystemModule)
	onSetCode = NewDefaultOnSetCode(systemModule)
}

func Test_DefaultOnSetCode_New(t *testing.T) {
	setupDefaultOnSetCode()

	expected := defaultOnSetCode{module: systemModule}

	assert.Equal(t, expected, onSetCode)
}

func Test_DefaultOnSetCode_SetCode(t *testing.T) {
	setupDefaultOnSetCode()

	codeBlob := sc.BytesToSequenceU8([]byte{1, 2, 3})

	systemModule.On("StorageCodeSet", codeBlob).Return()
	systemModule.On("DepositLog", primitives.NewDigestItemRuntimeEnvironmentUpgrade()).Return()
	systemModule.On("GetIndex").Return(sc.U8(moduleId))
	systemModule.On("DepositEvent", primitives.NewEvent(moduleId, EventCodeUpdated)).Return()

	err := onSetCode.SetCode(codeBlob)

	assert.Nil(t, err)
	systemModule.AssertCalled(t, "StorageCodeSet", codeBlob)
	systemModule.AssertCalled(t, "DepositLog", primitives.NewDigestItemRuntimeEnvironmentUpgrade())
	systemModule.AssertCalled(t, "GetIndex")
	systemModule.AssertCalled(t, "DepositEvent", primitives.NewEvent(sc.U8(moduleId), EventCodeUpdated))
}
