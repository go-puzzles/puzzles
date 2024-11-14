package discover

import (
	"sync"

	"github.com/go-puzzles/puzzles/cores/discover/manual"
)

type Service struct {
	ServiceName string
	Address     string
	Tags        []string
}

type ServiceFinder interface {
	GetAddress(service string) string
	GetAllAddress(service string) []string
	GetAddressWithTag(service, tag string) string
	GetAllAddressWithTag(service, tag string) []string

	RegisterService(service, address string) error
	RegisterServiceWithTag(service, address, tag string) error
	RegisterServiceWithTags(service, address string, tags []string) error
	Close()
}

var (
	defaultServiceFinder ServiceFinder
	finderMutex          sync.RWMutex
)

func init() {
	defaultServiceFinder = manual.NewDirectFinder()
}

func GetServiceFinder() ServiceFinder {
	finderMutex.RLock()
	defer finderMutex.RUnlock()
	return defaultServiceFinder
}

func SetFinder(finder ServiceFinder) {
	finderMutex.Lock()
	defer finderMutex.Unlock()
	defaultServiceFinder = finder
}

func GetAddress(srv string) string {
	return GetServiceFinder().GetAddress(srv)
}

func GetAddresses(srv string) []string {
	return GetServiceFinder().GetAllAddress(srv)
}

func GetAddressWithTag(srv, tag string) string {
	return GetServiceFinder().GetAddressWithTag(srv, tag)
}
