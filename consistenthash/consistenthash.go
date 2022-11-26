package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// hash 函数
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int // sorted
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 排序hash环
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	h := int(m.hash([]byte(key)))

	// 二分法
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= h
	})
	return m.hashMap[m.keys[index%len(m.keys)]]
}
