package main

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Sub[T any] struct {
	Id       uuid.UUID
	sendChan chan<- T
}

type LiveEventManager[T any] struct {
	subs      map[uuid.UUID]*Sub[T]
	subsMutex *sync.RWMutex
	buf       chan T
}

func newLiveEventManager[T any]() *LiveEventManager[T] {
	return &LiveEventManager[T]{
		subs:      make(map[uuid.UUID]*Sub[T]),
		subsMutex: &sync.RWMutex{},
		buf:       make(chan T),
	}
}

func (manager *LiveEventManager[T]) startBackgroundProcessing() {
	for msg := range manager.buf {
		// fmt.Printf("received message to broadcast: %+v\n", msg)
		// TODO: error recovery.
		manager.subsMutex.RLock()
		for _, sub := range manager.subs {
			sub.sendChan <- msg
		}
		manager.subsMutex.RUnlock()
	}
}

func (manager *LiveEventManager[T]) register(sendChan chan<- T) (*Sub[T], error) {
	if sendChan == nil {
		return nil, errors.New("sendChan cannot be nil")

	}
	sub := &Sub[T]{Id: uuid.New(), sendChan: sendChan}
	manager.subsMutex.Lock()
	defer manager.subsMutex.Unlock()
	manager.subs[sub.Id] = sub
	return sub, nil
}

func (manager *LiveEventManager[T]) deregister(sub *Sub[T]) error {
	if sub == nil {
		return errors.New("sub cannot be nil")
	}
	manager.subsMutex.Lock()
	defer manager.subsMutex.Unlock()

	delete(manager.subs, sub.Id)
	return nil
}

func (manager *LiveEventManager[T]) broadcast(msg T) {
	manager.buf <- msg
}
