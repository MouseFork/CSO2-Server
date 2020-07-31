package lock

import "sync"

var (
	GlobalMutex      sync.Mutex
	MainServerMutex  sync.Mutex
	UserManagerMutex sync.Mutex
	DataBaseMutex    sync.Mutex
)
