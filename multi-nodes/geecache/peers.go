package geecache

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	// PickPeer picks a peer according to the key.
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer. 对应HTTP客户端
type PeerGetter interface {
	// Get returns the value for the specified key and group.
	Get(in)
}
