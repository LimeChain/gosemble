//go:build nonwasmenv

package log

const (
	Critical = iota
	Warn
	Info
	Debug
	Trace
)

func Log(level int32, target []byte, message []byte) {
	panic(message)
}
