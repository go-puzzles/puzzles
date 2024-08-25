package discover

import (
	"sync"
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
	defaultServiceFinder = NewDirectFinder()
}

func GetServiceFinder() ServiceFinder {
	finderMutex.RLock()
	defer finderMutex.RUnlock()
	return defaultServiceFinder
}

func SetConsulFinderToDefault() {
	finderMutex.Lock()
	defer finderMutex.Unlock()
	defaultServiceFinder = GetConsulServiceFinder()
}

func GetConsulServiceFinder() ServiceFinder {
	return GetConsulClient()
}
