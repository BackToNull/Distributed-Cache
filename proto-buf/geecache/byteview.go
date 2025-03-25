package geecache

type ByteView struct {
	b []byte
}

// Len 实现Len方法，可以适配缓冲区的元素类型
func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// Notes: cloneBytes方法用于复制[]byte，防止缓存值被外部程序修改
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
