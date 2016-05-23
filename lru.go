//最少使用页面置换算法
package lru

import "container/list"

type T interface{}

type entity struct {
	key T
	val T
}

type Cache struct {
	maxEntity int
	ll        *list.List
	m         map[T]*list.Element
	hook      func(T, T) error
}

func New(maxEntity int, hook func(T, T) error) *Cache {
	return &Cache{
		maxEntity: maxEntity,
		ll:        list.New(),
		m:         make(map[T]*list.Element),
		hook:      hook,
	}
}

func (c *Cache) Add(key, val T) {
	if c.ll == nil {
		c.ll = list.New()
		c.m = make(map[T]*list.Element)
	}
	if ele, ok := c.m[key]; ok {
		ele.Value.(*entity).val = val
		c.ll.MoveToFront(ele)
		return
	}
	c.m[key] = c.ll.PushFront(&entity{key: key, val: val})

	if c.maxEntity > 0 && c.maxEntity <= c.Len() {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key T) (val T, ok bool) {
	if c.ll == nil {
		return
	}
	if ele, ok := c.m[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entity).val, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	if c.ll == nil {
		return
	}
	c.removeElement(c.ll.Back())
}

func (c *Cache) Remove(key T) {
	if c.ll == nil {
		return
	}
	if ele, ok := c.m[key]; ok {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entity)
	delete(c.m, kv.key)
	if c.hook != nil {
		c.hook(kv.key, kv.val)
	}
}

func (c *Cache) Len() int {
	if c.ll == nil {
		return 0
	}
	return c.ll.Len()
}
