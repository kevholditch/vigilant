package utils

import (
	"sort"
	"sync"
)

// OrderedMap maintains a map for fast lookups and a slice for consistent ordering
// This version is specifically for string keys to enable proper sorting
type OrderedMap[V any] struct {
	items map[string]V
	order []string
	mutex sync.RWMutex
}

// NewOrderedMap creates a new OrderedMap
func NewOrderedMap[V any]() *OrderedMap[V] {
	return &OrderedMap[V]{
		items: make(map[string]V),
		order: make([]string, 0),
	}
}

// Set adds or updates a key-value pair and maintains order
func (om *OrderedMap[V]) Set(key string, value V) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Check if key already exists
	_, exists := om.items[key]

	// Add/update the item
	om.items[key] = value

	// Add to order if it's a new key
	if !exists {
		om.order = append(om.order, key)
		om.sortOrder()
	}
}

// Get retrieves a value by key
func (om *OrderedMap[V]) Get(key string) (V, bool) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	value, exists := om.items[key]
	return value, exists
}

// Delete removes a key-value pair
func (om *OrderedMap[V]) Delete(key string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	delete(om.items, key)
	om.removeFromOrder(key)
}

// Keys returns all keys in order
func (om *OrderedMap[V]) Keys() []string {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	result := make([]string, len(om.order))
	copy(result, om.order)
	return result
}

// Values returns all values in order
func (om *OrderedMap[V]) Values() []V {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	result := make([]V, 0, len(om.items))
	for _, key := range om.order {
		if value, exists := om.items[key]; exists {
			result = append(result, value)
		}
	}
	return result
}

// Len returns the number of items
func (om *OrderedMap[V]) Len() int {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	return len(om.items)
}

// Clear removes all items
func (om *OrderedMap[V]) Clear() {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	om.items = make(map[string]V)
	om.order = make([]string, 0)
}

// sortOrder sorts the order slice
func (om *OrderedMap[V]) sortOrder() {
	sort.Strings(om.order)
}

// removeFromOrder removes a key from the order slice
func (om *OrderedMap[V]) removeFromOrder(key string) {
	for i, k := range om.order {
		if k == key {
			om.order = append(om.order[:i], om.order[i+1:]...)
			break
		}
	}
}
