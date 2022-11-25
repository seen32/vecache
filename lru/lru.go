package lru

import (
	"container/list"
)

// 缓存值的抽象接口
type Value interface {
	// 缓存值所占用的内存大小
	Len() int
}

// 双向链表的节点
type entry struct {
	key   string
	value Value
}

type Cache struct {
	// 允许使用的最大内存
	maxBytes int64
	// 当前已使用的内存
	nbytes int64

	// double linked list
	dlist *list.List
	// hashmap
	cache map[string]*list.Element

	// 当缓存被删除时调用的回调函数
	onEvicted func(key string, value Value)
}

// 初始化函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nbytes:    0,
		dlist:     list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Len() int {
	return c.dlist.Len()
}

// 获取缓存的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if el, ok := c.cache[key]; ok {
		// 最近访问的缓存放到dlist的表头
		c.dlist.MoveToFront(el)
		kv := el.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Set(key string, value Value) {
	if el, ok := c.cache[key]; ok {
		c.dlist.MoveToFront(el)
		kv := el.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		el := c.dlist.PushFront(&entry{key, value})
		c.cache[key] = el
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 如果内存超出，则触发缓存淘汰
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 删除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.dlist.Back()
	if ele != nil {
		c.dlist.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		// 调用回调函数
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}
