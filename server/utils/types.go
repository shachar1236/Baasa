package utils

import (
	"sync"
)

type ValueWithMutex[T any] struct {
    Value T
    m sync.Mutex
}

func (this *ValueWithMutex[T]) Lock() {
    this.m.Lock()
}

func (this *ValueWithMutex[T]) Unlock() {
    this.m.Unlock()
}

func (this *ValueWithMutex[T]) TryLock() bool {
    return this.m.TryLock()
}
