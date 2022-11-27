package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"vecache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *vecache.Group {
	return vecache.NewGroup("scores", 2<<10, vecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, ve *vecache.Group) {
	peers := vecache.NewHTTPPool(addr[7:])
	peers.Set(addrs...)
	ve.RegisterPeers(peers)
	peers.Start()
}

// 开启网关
func startAPIServer(apiAddr string, ve *vecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := ve.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int // 命令行参数：端口号
	var api bool // 命令行参数：开启网关
	flag.IntVar(&port, "port", 8001, "vecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	// 网关IP
	apiAddr := "http://localhost:8888"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	ve := createGroup()
	if api {
		go startAPIServer(apiAddr, ve)
	}
	// 开启缓存服务
	startCacheServer(addrMap[port], addrs, ve)
}
