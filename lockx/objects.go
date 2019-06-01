package lockx

import "sync"

type ResourceLock struct {
	sync.Mutex
	tryLock sync.Mutex
	locked  bool
	action  Action
}

func (rl *ResourceLock) TryLock(action Action) (bool, Action) {
	if rl.locked {
		return false, rl.action
	}
	rl.tryLock.Lock()
	defer rl.tryLock.Unlock()
	if rl.locked {
		return false, rl.action
	}
	rl.Mutex.Lock()
	rl.locked = true
	rl.action = action
	return true, action
}

func (rl *ResourceLock) Unlock() {
	rl.Mutex.Unlock()
	rl.locked = false
	rl.action = nil
}

func (rl *ResourceLock) IsLocked() bool {
	return rl.locked
}

type Action interface {
	GetCode() string
	GetMessage() string
	GetOwnerKey() string
}
