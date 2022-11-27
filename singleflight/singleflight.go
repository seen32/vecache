package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 防止缓存击穿，保证并发请求时，只向目标远程节点发送一次请求
type Group struct {
	mutex sync.Mutex
	maps  map[string]*call
}

// 借鉴了 internal/singleflight 包下的 同名方法
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mutex.Lock()
	if g.maps == nil {
		g.maps = make(map[string]*call)
	}
	// 如果正在进行处理，则等待
	if c, ok := g.maps[key]; ok {
		g.mutex.Unlock()
		// wait等待处理完
		c.wg.Wait()
		// 此时已经处理完，返回的结果已经存放到call结构中
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.maps[key] = c

	g.mutex.Unlock()
	// 处理...
	c.val, c.err = fn()

	c.wg.Done()

	g.mutex.Lock()
	delete(g.maps, key)
	g.mutex.Unlock()

	return c.val, c.err
}
