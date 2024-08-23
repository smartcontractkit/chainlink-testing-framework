package seth

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type cacheItem struct {
	header    *types.Header
	frequency int
}

// LFUHeaderCache is a Least Frequently Used header cache
type LFUHeaderCache struct {
	capacity uint64
	mu       *sync.RWMutex
	cache    map[int64]*cacheItem //key is block number
}

// NewLFUBlockCache creates a new LFU cache with the given capacity.
func NewLFUBlockCache(capacity uint64) *LFUHeaderCache {
	return &LFUHeaderCache{
		capacity: capacity,
		cache:    make(map[int64]*cacheItem),
		mu:       &sync.RWMutex{},
	}
}

// Get retrieves a header from the cache.
func (c *LFUHeaderCache) Get(blockNumber int64) (*types.Header, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.cache[blockNumber]; found {
		item.frequency++
		L.Trace().Msgf("Found header %d in cache", blockNumber)
		return item.header, true
	}
	return nil, false
}

// Set adds or updates a header in the cache.
func (c *LFUHeaderCache) Set(header *types.Header) error {
	if header == nil {
		return fmt.Errorf("header is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if oldHeader, found := c.cache[header.Number.Int64()]; found {
		L.Trace().Msgf("Setting header %d in cache", header.Number.Int64())
		c.cache[int64(header.Number.Int64())] = &cacheItem{header: header, frequency: oldHeader.frequency + 1}
		return nil
	}

	if uint64(len(c.cache)) >= c.capacity {
		c.evict()
	}
	L.Trace().Msgf("Setting header %d in cache", header.Number.Int64())
	c.cache[int64(header.Number.Int64())] = &cacheItem{header: header, frequency: 1}

	return nil
}

// evict removes the least frequently used item from the cache. If more than one item has the same frequency, the oldest is evicted.
func (c *LFUHeaderCache) evict() {
	var leastFreq int = int(^uint(0) >> 1)
	var evictKey int64
	oldestBlockNumber := ^uint64(0)
	for key, item := range c.cache {
		if item.frequency < leastFreq {
			evictKey = key
			leastFreq = item.frequency
			oldestBlockNumber = item.header.Number.Uint64()
		} else if item.frequency == leastFreq && item.header.Number.Uint64() < oldestBlockNumber {
			// If frequencies are the same, evict the oldest based on block number
			evictKey = key
			oldestBlockNumber = item.header.Number.Uint64()
		}
	}
	L.Trace().Msgf("Evicted header %d from cache", evictKey)
	delete(c.cache, evictKey)
}
