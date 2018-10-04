//+build windows

package log

import (
	"github.com/inconshreveable/log15"
)

type SysLogHandler struct {
}

func NewSyslog(sec *ini.Section, format log15.Format) *SysLogHandler {
	return &SysLogHandler{}
}

func (sw *SysLogHandler) Log(r *log15.Record) error {
	return nil
}

func (sw *SysLogHandler) Close() {
}
