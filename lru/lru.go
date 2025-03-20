package lru

import "container/list"

type Cache struct {
	// 缓存最大容量
	maxBytes int64
	// 当前已使用的内存
	nbytes int64
	//双向链表
	//back: 待删除的节点
	//front: 最近访问的节点
	ll *list.List
	// Element is an element of a linked list.
	cache map[string]*list.Element
	// option and executed when an entry is purged 淘汰队首节点时，需要用key从字典删除对应的映射
	OnEvicted func(key string, value Value)
}

// entry 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value 允许值是实现了Value接口的任意类型,保证通用性
type Value interface {
	Len() int
}

// New Cache的构造函数
func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//update kv pair
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//insert kv pair
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//如果超过了设定的最大内存，则移除最少访问的节点
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
