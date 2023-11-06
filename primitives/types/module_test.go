package types

import (
	"strconv"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

const (
	numModules = 4
)

var (
	modules []Module
)

type module struct {
	Module
	Index sc.U8
}

func newTestModule(index sc.U8) Module {
	return module{
		Index: index,
	}
}

func (m module) GetIndex() sc.U8 {
	return m.Index
}

func setup() {
	modules = []Module{newTestModule(sc.U8(0)), newTestModule(sc.U8(1)), newTestModule(sc.U8(2)), newTestModule(sc.U8(3))}
}

func Test_GetModule(t *testing.T) {
	setup()
	for i := 0; i < numModules; i++ {
		m, err := GetModule(sc.U8(i), modules)
		assert.NoError(t, err)
		assert.Equal(t, m.GetIndex(), sc.U8(i))
		assert.Equal(t, m, modules[i])
	}
}

func Test_GetModule_FailWhenNonExistent(t *testing.T) {
	setup()
	m, err := GetModule(sc.U8(numModules), modules)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func Test_MustGetModule(t *testing.T) {
	setup()
	for i := 0; i < numModules; i++ {
		m, err := GetModule(sc.U8(i), modules)

		assert.NoError(t, err)
		assert.Equal(t, m.GetIndex(), sc.U8(i))
		assert.Equal(t, m, modules[i])
	}
}

func Test_MustGetModule_ErrorNonExistent(t *testing.T) {
	setup()

	mod, err := GetModule(sc.U8(numModules), modules)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "module with index ["+strconv.Itoa(int(numModules))+"] not found.")
	assert.Nil(t, mod)
}
