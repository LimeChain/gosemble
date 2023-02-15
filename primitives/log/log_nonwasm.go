//go:build nonwasmenv

package log

func Log(level int32, target []byte, message []byte) {
	panic("not implemented")
}
