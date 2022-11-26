package main

import (
	"fmt"
	"vecache"
)

func main() {
	db := map[string]string{
		"ZhangSan": "176",
		"LiSi":     "172",
		"WangWu":   "184",
	}

	// 创建group
	vecache.NewGroup("height", 2<<10, vecache.GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			return []byte(v), nil
		} else {
			return []byte{}, fmt.Errorf("Can not Get From DataSource")
		}
	}))

	// 启动HTTP服务
	addr := "localhost:8080"
	pool := vecache.NewHTTPPool(addr)
	pool.Start()

	// 通过命令行curl进行测试
}
