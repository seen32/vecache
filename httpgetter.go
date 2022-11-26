package vecache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// http服务客户端，从其他的缓存节点获取缓存
type HTTPGetter struct {
	baseURL string
}

func (h *HTTPGetter) Get(group string, key string) ([]byte, error) {
	url := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

var _ PeerGetter = (*HTTPGetter)(nil)
