package main

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/radkomih/gosemble/helpers"
	"github.com/stretchr/testify/require"
)

func Test_Core_version(t *testing.T) {
	rt := helpers.NewTestInstanceWithTrie(t, &trie.Trie{})

	res, err := rt.Exec("Core_version", []byte{})

	fmt.Println(res)

	require.NoError(t, err)
}
