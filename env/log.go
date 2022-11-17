package env

/*
	Log: Request to print a log message on the host. Note that this will be
	only displayed if the host is enabled to display log messages with given level and target.
*/
//go:wasm-module env
//go:export ext_logging_log_version_1
func extLoggingLogVersion1(level int32, target int64, message int64)

func ExtLoggingLogVersion1(level int32, target int64, message int64) {
	extLoggingLogVersion1(level, target, message)
}
