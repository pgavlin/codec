package typecache

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

// Eventually consistent cache mapping Go types to values.
type Cache[T any] struct {
	// Note: using a uintptr as key instead of reflect.Type shaved ~15ns off of
	// the ~30ns Marhsal/Unmarshal functions which were dominated by the map
	// lookup time for simple types like bool, int, etc..
	cache unsafe.Pointer // map[unsafe.Pointer]codec
}

func (c Cache[T]) load() map[unsafe.Pointer]T {
	p := atomic.LoadPointer(&c.cache)
	return *(*map[unsafe.Pointer]T)(unsafe.Pointer(&p))
}

func (c Cache[T]) set(type_ reflect.Type, val T, oldCache map[unsafe.Pointer]T) {
	newCache := make(map[unsafe.Pointer]T, len(oldCache)+1)
	newCache[typeid(type_)] = val

	for t, c := range oldCache {
		newCache[t] = c
	}

	atomic.StorePointer(&c.cache, *(*unsafe.Pointer)(unsafe.Pointer(&newCache)))
}

// Gets the value associated with the given type. If the type is not in the cache,
// calls create to create a new value and records it in the cache.
func (c Cache[T]) GetOrCreate(type_ reflect.Type, create func(reflect.Type) T) T {
	cache := c.load()
	v, found := cache[typeid(type_)]
	if !found {
		v = create(type_)
		c.set(type_, v, cache)
	}
	return v
}

func typeid(t reflect.Type) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&t)).ptr
}

type iface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}
