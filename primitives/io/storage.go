package io

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Storage interface {
	Append(key []byte, value []byte)
	Clear(key []byte)
	ClearPrefix(key []byte, limit []byte)
	Exists(key []byte) bool
	Get(key []byte) sc.Option[sc.Sequence[sc.U8]]
	NextKey(key int64) int64
	Read(key []byte, valueOut []byte, offset int32) sc.Option[sc.U32]
	Root(version int32) []byte
	Set(key []byte, value []byte)
}

type storage struct {
}

func NewStorage() Storage {
	return storage{}
}

func (s storage) Append(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageAppendVersion1(keyOffsetSize, valueOffsetSize)
}

func (s storage) Clear(key []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	env.ExtStorageClearVersion1(keyOffsetSize)
}

func (s storage) ClearPrefix(key []byte, limit []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	limitOffsetSize := utils.BytesToOffsetAndSize(limit)
	env.ExtStorageClearPrefixVersion2(keyOffsetSize, limitOffsetSize)
}

func (s storage) Exists(key []byte) bool {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	return env.ExtStorageExistsVersion1(keyOffsetSize) != 0
}

func (s storage) Get(key []byte) sc.Option[sc.Sequence[sc.U8]] {
	value := get(key)

	buffer := &bytes.Buffer{}
	buffer.Write(value)

	return sc.DecodeOption[sc.Sequence[sc.U8]](buffer)
}

func (s storage) NextKey(key int64) int64 {
	panic("not implemented")
}

func (s storage) Read(key []byte, valueOut []byte, offset int32) sc.Option[sc.U32] {
	value := read(key, valueOut, offset)

	buffer := &bytes.Buffer{}
	buffer.Write(value)

	return sc.DecodeOption[sc.U32](buffer)
}

func (s storage) Root(version int32) []byte {
	valueOffsetSize := env.ExtStorageRootVersion2(version)
	offset, size := utils.Int64ToOffsetAndSize(valueOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

func (s storage) Set(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageSetVersion1(keyOffsetSize, valueOffsetSize)
}

// get gets the value from storage by the provided key. The wasm memory slice (value)
// represents an encoded Option<sc.Sequence[sc.U8]> (option of encoded slice).
func get(key []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := env.ExtStorageGetVersion1(keyOffsetSize)
	offset, size := utils.Int64ToOffsetAndSize(valueOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

// read reads the given key value from storage, placing the value into buffer valueOut depending on offset.
// The wasm memory slice represents an encoded Option<sc.U32> representing the number of bytes left at supplied offset.
func read(key []byte, valueOut []byte, offset int32) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOutOffsetSize := utils.BytesToOffsetAndSize(valueOut)

	resultOffsetSize := env.ExtStorageReadVersion1(keyOffsetSize, valueOutOffsetSize, offset)
	offset, size := utils.Int64ToOffsetAndSize(resultOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)

	return value
}

type TransactionBroker interface {
	Start()
	Commit()
	Rollback()
}

type transactionBroker struct{}

func NewTransactionBroker() TransactionBroker {
	return transactionBroker{}
}

// Start a new nested transaction.
//
// This allows to either commit or roll back all changes that are made after this call.
// For every transaction there must be a matching call to either `rollback_transaction`
// or `commit_transaction`. This is also effective for all values manipulated using the
// `DefaultChildStorage` API.
//
// # Warning
//
// This is a low level API that is potentially dangerous as it can easily result
// in unbalanced transactions. For example, FRAME users should use high level storage
// abstractions.
func (tb transactionBroker) Start() {
	env.ExtStorageStartTransactionVersion1()
}

// Rollback the last transaction started by `start_transaction`.
//
// Any changes made during that transaction are discarded.
//
// # Panics
//
// Will panic if there is no open transaction.
func (tb transactionBroker) Rollback() {
	env.ExtStorageRollbackTransactionVersion1() // TODO: .expect("No open transaction that can be rolled back.");
}

// Commit the last transaction started by `start_transaction`.
//
// Any changes made during that transaction are committed.
//
// # Panics
//
// Will panic if there is no open transaction.
func (tb transactionBroker) Commit() {
	env.ExtStorageCommitTransactionVersion() // TODO: .expect("No open transaction that can be committed.");
}
