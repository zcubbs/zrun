// Package healthcheck
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package healthcheck

import (
	"github.com/zcubbs/zrun/log"
	"net"
	"strings"
	"sync"
	"time"
)

// Services is a string array storing
// the services that are to be waited for
type Services []string

// Set is used to append a string
// to the service, to implement
// the interface flag.Value
func (s *Services) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// String returns a string
// representation of the flag,
// to implement the interface
// flag.Value
func (s *Services) String() string {
	return strings.Join(*s, ", ")
}

// Wait waits for all services
func Wait(services []string, tSeconds int) bool {
	t := time.Duration(tSeconds) * time.Second
	now := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(services))

	success := make(chan bool, 1)

	go func() {
		for _, service := range services {
			go waitOne(service, &wg, now)
		}
		wg.Wait()
		success <- true
	}()

	select {
	case <-success:
		return true
	case <-time.After(t):
		return false
	}
}

func waitOne(service string, wg *sync.WaitGroup, start time.Time) {
	defer wg.Done()
	for {
		_, err := net.Dial("tcp", service)
		if err == nil {
			log.Fatalw("%s is available after %s", map[string]any{
				"service": service,
				"elapsed": time.Since(start),
			})
			break
		}
		time.Sleep(time.Second)
	}
}
