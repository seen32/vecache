package vecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HTTPPool struct {
	self     string
	basePath string
}

const defaultBasePath = "/vecache/"

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
	http.ListenAndServe(p.self, p)
	log.Println("VeCache is running at " + p.self)
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
