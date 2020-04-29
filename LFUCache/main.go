package main

import (
	"container/list"
	"fmt"
)

type LFUData struct {
	key   int
	value int
	count int
}

type LFUCache struct {
	capacity int
	findMap  map[int]*list.Element
	countMap map[int]*list.List
	minCount int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		capacity: capacity,
		findMap:  make(map[int]*list.Element, capacity),
		countMap: make(map[int]*list.List, capacity),
	}
}

func (l *LFUCache) Get(key int) int {
	d, ok := l.findMap[key]
	if !ok {
		return -1
	}
	return l.update(d)
}

func (l *LFUCache) update(ele *list.Element) int {
	data := ele.Value.(*LFUData)
	curList, ok := l.countMap[data.count]
	if !ok {
		return -1
	}
	curList.Remove(ele)
	if curList.Len() == 0 {
		if data.count == l.minCount {
			l.minCount++
		}
	}
	data.count++
	newList, ok := l.countMap[data.count]
	if !ok {
		newList = list.New()
		l.countMap[data.count] = newList
	}
	newE := newList.PushBack(data)
	l.findMap[data.key] = newE
	return data.value
}

func (l *LFUCache) Put(key int, value int) {
	existEle, ok := l.findMap[key]
	if ok {
		data := existEle.Value.(*LFUData)
		data.value = value
		l.update(existEle)
		return
	}
	if len(l.findMap) == l.capacity {
		minList, ok := l.countMap[l.minCount]
		if !ok {
			return
		}
		delE := minList.Front()
		minList.Remove(delE)
		data := delE.Value.(*LFUData)
		delete(l.findMap, data.key)
	}
	newData := &LFUData{
		key:   key,
		value: value,
		count: 1,
	}
	newList, ok := l.countMap[newData.count]
	if !ok {
		newList = list.New()
		l.countMap[newData.count] = newList
	}
	newEle := newList.PushBack(newData)
	l.findMap[key] = newEle
	l.minCount = newData.count
}

func main() {
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // 返回 1
	cache.Put(3, 3)           // 去除 key 2
	fmt.Println(cache.Get(2)) // 返回 -1 (未找到key 2)
	fmt.Println(cache.Get(3)) // 返回 3
	cache.Put(4, 4)           // 去除 key 1
	fmt.Println(cache.Get(1)) // 返回 -1 (未找到 key 1)
	fmt.Println(cache.Get(3)) // 返回 3
	fmt.Println(cache.Get(4)) // 返回 4

}
