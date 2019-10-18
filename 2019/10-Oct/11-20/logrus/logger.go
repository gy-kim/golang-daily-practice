package logrus

import (
	"io"
	"sync"
)

type Logger struct {
	Out          io.Writer
	Hooks        LevelHooks
	Formatter    Formatter
	ReportCaller bool
	Level        Level
	mu           MutexWrap
	entryPool    sync.Pool
	ExitFunc     exitFunc
}

type exitFunc func(int)

type MutexWap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWap) Disable() {
	mw.disabled = true
}
