package proxy

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	backendFilterAll       = iota
	backendFilterFailed    = iota
	backendFilterAvailable = iota
)

var ErrAlreadyUpdateStarted = errors.New("Periodical updating is already running")

type BackendUpdateFunc func() ([]*Backend, error)

type Router struct {
	sync.RWMutex
	backends       []*Backend
	UpdateFunc     BackendUpdateFunc
	UpdateInterval time.Duration
	updateStopCh   chan bool
}

func NewRouter(interval time.Duration, updateFunc BackendUpdateFunc) *Router {
	return &Router{
		backends:       []*Backend{},
		UpdateFunc:     updateFunc,
		UpdateInterval: interval,
		updateStopCh:   nil,
	}
}

func (router *Router) StartUpdate() error {
	router.Lock()
	defer router.Unlock()

	if router.updateStopCh != nil {
		return ErrAlreadyUpdateStarted
	}

	router.updateStopCh = make(chan bool)

	go func() {
		for {
			select {
			case <-router.updateStopCh:
				return
			case <-time.After(router.UpdateInterval):
				router.Update()
			}
		}
	}()

	log.Println("Started periodical update of backends")

	return nil
}

func (router *Router) StopUpdate() {
	router.Lock()
	defer router.Unlock()

	if router.updateStopCh != nil {
		router.updateStopCh <- true
		log.Println("Stopped periodical update of backends")
	}
}

func (router *Router) Update() {
	router.Lock()
	defer router.Unlock()

	newBackends, err := router.UpdateFunc()
	if err == nil {
		router.backends = newBackends
	} else {
		log.Printf("Failed to update backends: %s", err.Error())
	}
}

func (router *Router) getBackends(filter int) []*Backend {
	router.RLock()
	defer router.RUnlock()

	filteredBackends := make([]*Backend, 0, len(router.backends))

	for _, backend := range router.backends {
		switch filter {
		case backendFilterAll:
		case backendFilterFailed:
			if backend.Available {
				continue
			}
		case backendFilterAvailable:
			if !backend.Available {
				continue
			}
		}

		filteredBackends = append(filteredBackends, backend)
	}

	return filteredBackends
}

func (router *Router) Backends() []*Backend {
	return router.getBackends(backendFilterAll)
}

func (router *Router) FailedBackends() []*Backend {
	return router.getBackends(backendFilterFailed)
}

func (router *Router) AvailableBackends() []*Backend {
	return router.getBackends(backendFilterAvailable)
}

func (router *Router) ShuffledAvailableBackends() []*Backend {
	backends := router.AvailableBackends()
	shuffledBackends := make([]*Backend, len(backends))

	pattern := rand.Perm(len(backends))
	for i, idx := range pattern {
		shuffledBackends[i] = backends[idx]
	}

	return shuffledBackends
}
