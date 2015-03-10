package proxy

import (
	"log"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	sync.Mutex
	Url               *url.URL
	Available         bool
	nextCheckInterval time.Duration
	resumeTimer       *time.Timer
}

func NewBackend(url *url.URL) *Backend {
	return &Backend{
		Url:               url,
		Available:         true,
		nextCheckInterval: time.Duration(time.Second) * 15,
		resumeTimer:       nil,
	}
}

func (backend *Backend) Fail() {
	backend.Lock()
	defer backend.Unlock()

	if !backend.Available {
		return
	}

	backend.Available = false
	checkInterval := backend.nextCheckInterval
	backend.nextCheckInterval = checkInterval * 2

	backend.resumeTimer = time.AfterFunc(backend.nextCheckInterval, func() {
		backend.Lock()
		backend.Available = true
		backend.Unlock()

		log.Printf("Backend %s resumed (automatically)", backend.Url.String())
	})

	log.Printf("Backend %s marked as failure, will resume after %s", backend.Url.String(), checkInterval.String())
}

func (backend *Backend) Ok() {
	backend.Lock()
	defer backend.Unlock()

	backend.Available = true
	backend.nextCheckInterval = time.Duration(time.Second) * 15

	if backend.resumeTimer != nil {
		backend.resumeTimer.Stop()
		backend.resumeTimer = nil
	}

	log.Printf("Backend %s resumed", backend.Url.String())
}
