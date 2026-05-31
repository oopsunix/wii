package provider

import (
	"context"
	"sync"
)

// Provider defines the interface for a package manager integration.
type Provider interface {
	// Name returns the provider name (e.g. "brew", "cargo").
	Name() string

	// Available reports whether the package manager is installed on this system.
	Available() bool

	// Fetch queries installed packages and populates the cache.
	// Implementations must respect ctx for cancellation and timeouts.
	Fetch(ctx context.Context) error
}

// versionCache is a thread-safe cache for command name -> version mappings.
var versionCache sync.Map

// CacheGet retrieves a cached version for the given command name.
func CacheGet(name string) (string, bool) {
	v, ok := versionCache.Load(name)
	if !ok {
		return "", false
	}
	return v.(string), true
}

// CacheSet stores a version for the given command name.
func CacheSet(name, version string) {
	versionCache.Store(name, version)
}

// CacheKeys returns all cached command names.
func CacheKeys() []string {
	var keys []string
	versionCache.Range(func(key, _ any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

var registry []Provider

// Register adds a provider to the global registry.
// Providers self-register via init() in their respective files.
func Register(p Provider) {
	registry = append(registry, p)
}

// Resolve returns all providers that are available on the current system.
func Resolve() []Provider {
	var available []Provider
	for _, p := range registry {
		if p.Available() {
			available = append(available, p)
		}
	}
	return available
}

// FetchAll queries all available providers in parallel.
func FetchAll(ctx context.Context) {
	providers := Resolve()
	if len(providers) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // limit concurrent provider queries

	for _, p := range providers {
		wg.Add(1)
		sem <- struct{}{}
		go func(p Provider) {
			defer wg.Done()
			defer func() { <-sem }()
			_ = p.Fetch(ctx) // best-effort; failures are silently ignored
		}(p)
	}
	wg.Wait()
}
