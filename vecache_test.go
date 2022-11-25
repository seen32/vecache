package vecache

import (
	"fmt"
	"testing"
)

func TestGroup_Get(t *testing.T) {
	// 模拟数据库
	db := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	group := NewGroup("test", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			return []byte(v), nil
		} else {
			return []byte{}, fmt.Errorf("Can not Get DataSource")
		}
	}))

	for k, v := range db {
		view, err := group.Get(k)
		if err != nil {
			t.Fatalf("cache Get error")
		}
		if view.String() != v {
			t.Fatalf("cache miss")
		}
	}
}
