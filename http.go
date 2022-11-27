package vecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"vecache/consistenthash"
)

type HTTPPool struct {
	self        string
	basePath    string
	mutex       sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*HTTPGetter
}

const (
	defaultBasePath = "/vecache/"
	defaultReplicas = 50
)

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) Start() {
	log.Println("VeCache is running at " + p.self)
	log.Fatal(http.ListenAndServe(p.self, p))
}

func (p *HTTPPool) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if !strings.HasPrefix(path, p.basePath) {
		panic("URL unresolved : " + path)
	}
	p.Log("%s %s", request.Method, path)

	// 格式
	// URL Format : /[basePath]/[groupName]/[key]

	// 将bashPath后的部分根据 "/" 进行拆分成至多2份
	parts := strings.SplitN(path[len(p.basePath):], "/", 2)

	if len(parts) != 2 {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group, err := GetGroup(groupName)
	if err != nil {
		http.Error(writer, "Unknown Group", http.StatusNotFound)
		return
	}

	byteView, err := group.Get(key)
	if err != nil {
		http.Error(writer, "Unknown key", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Write(byteView.ByteSlice())
}

// 添加远程节点
func (p *HTTPPool) Set(peers ...string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)

	p.httpGetters = make(map[string]*HTTPGetter, len(peers))

	for _, peer := range peers {
		p.httpGetters[peer] = &HTTPGetter{baseURL: peer + p.basePath}
	}
}

// 选取远程节点
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer[7:] != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)
