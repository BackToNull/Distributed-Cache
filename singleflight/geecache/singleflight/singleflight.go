package singleflight

import "sync"

// singleflight 针对相同的key，8003节点并发地向8001发起三次请求是没有必要的，这个包将其优化成只向远端节点发起一次请求

//缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。

//缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。

//缓存穿透：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

// call 代表正在进行中，或已经结束的请求。使用WaitGroup锁来避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是singleflight的主数据结构，管理不同key的请求(call)
type Group struct {
	mu sync.Mutex //保护成员变量m不被并发读写而加上的锁
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 延迟初始化
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//如果有请求正在进行中，则等待
		c.wg.Wait()
		//请求结束，返回结果
		return c.val, c.err
	}

	c := new(call)
	//发起请求前加锁
	c.wg.Add(1)
	//添加到g.m，表明key已经有对应的请求在处理
	g.m[key] = c
	g.mu.Unlock()
	//调用fn，发起请求
	c.val, c.err = fn()
	//请求结束
	c.wg.Done()

	g.mu.Lock()
	//更新g.m
	delete(g.m, key)
	g.mu.Unlock()
	//返回结果
	return c.val, c.err
}
