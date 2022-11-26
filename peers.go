package vecache

// 根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 与HTTPPool相对应，PeerGetter是http客户端
type PeerGetter interface {
	// Get() 方法用于从对应 group 查找缓存值
	Get(group string, key string) ([]byte, error)
}
