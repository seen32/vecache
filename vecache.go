package vecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 当缓存不存在时，通过调用回调函数，从数据源获取数据并添加到缓存中
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 缓存命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[VeCache] hit")
		return v, nil
	}
	// 如果没有缓存，从数据源获取缓存
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	tmp := make([]byte, len(bytes))
	copy(tmp, bytes)
	value := ByteView{b: tmp}
	g.mainCache.Set(key, value)
	return value, nil
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mutex.Lock()
	defer mutex.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) (*Group, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if g, ok := groups[name]; ok {
		return g, nil
	} else {
		return nil, fmt.Errorf("Unknown Group")
	}
}
