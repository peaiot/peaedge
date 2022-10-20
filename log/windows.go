//go:build windows
// +build windows

package log

import "github.com/peaiot/logging"

func SetupSyslog(level logging.Level, syslogaddr string, module string) logging.LeveledBackend {
	return nil
}
