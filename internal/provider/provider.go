package provider

import (
	"context"
	"sync"

	"wii/internal/model"
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

// NamesProvider is optionally implemented by providers that can list
// their installed package names (for whitelist filtering).
type NamesProvider interface {
	Provider
	FetchNames(ctx context.Context) map[string]bool
}

// EntryProvider is optionally implemented by providers that generate
// package-level entries directly (e.g. Homebrew formulae instead of binaries).
type EntryProvider interface {
	NamesProvider
	FetchEntries(ctx context.Context) []model.Tool
}

// ResolveNames queries all available providers that support name listing
// and returns a map of provider label -> set of installed names.
func ResolveNames(ctx context.Context) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	for _, p := range Resolve() {
		if np, ok := p.(NamesProvider); ok {
			names := np.FetchNames(ctx)
			if len(names) > 0 {
				result[np.Name()] = names
			}
		}
	}
	return result
}

// ResolveEntries queries all available providers that support entry generation
// and returns a map of provider label -> entries.
func ResolveEntries(ctx context.Context) map[string][]model.Tool {
	result := make(map[string][]model.Tool)
	for _, p := range Resolve() {
		if ep, ok := p.(EntryProvider); ok {
			entries := ep.FetchEntries(ctx)
			if len(entries) > 0 {
				result[ep.Name()] = entries
			}
		}
	}
	return result
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
