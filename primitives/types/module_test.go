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
		m, ok := GetModule(sc.U8(i), modules)
		assert.True(t, ok)
		assert.Equal(t, m.GetIndex(), sc.U8(i))
		assert.Equal(t, m, modules[i])
	}
}

func Test_GetModule_FailWhenNonExistent(t *testing.T) {
	setup()
	m, ok := GetModule(sc.U8(numModules), modules)
	assert.False(t, ok)
	assert.Nil(t, m)
}

func Test_MustGetModule(t *testing.T) {
	setup()
	for i := 0; i < numModules; i++ {
		m := MustGetModule(sc.U8(i), modules)
		assert.Equal(t, m.GetIndex(), sc.U8(i))
		assert.Equal(t, m, modules[i])
	}
}

func Test_MustGetModule_PanicWhenNonExistent(t *testing.T) {
	setup()
	assert.PanicsWithValue(t, "module ["+strconv.Itoa(int(numModules))+"] not found.", func() {
		MustGetModule(sc.U8(numModules), modules)
	})
}
