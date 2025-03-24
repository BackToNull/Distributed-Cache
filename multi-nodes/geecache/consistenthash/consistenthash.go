package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash map bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash Hash
	// replicas is the number of virtual nodes for each real node
	replicas int
	// keys 哈希环
	keys []int
	// hashMap 虚拟节点与真实节点的映射表
	hashMap map[int]string
}

// New creates a Map instance
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

// Add add some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//虚拟节点哈希值的计算方法是：将真实节点的key和i拼接成新的字符串，然后对这个字符串计算哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	//保持哈希环的有序性
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//取余是因为如果idx == len(m.keys)，此时应该选择m.keys[0]
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
