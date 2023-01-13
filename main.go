package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const Test = "value"

type Cache struct {
	mutex           *sync.RWMutex
	records         map[string]*list.Element
	cleanupInterval time.Duration
	queue           *list.List
	capacity        int64
}

type Record struct {
	Value interface{}
	Key   string
}

func NewCache(cap int64) *Cache {
	return &Cache{
		mutex:    &sync.RWMutex{},
		records:  make(map[string]*list.Element, cap),
		capacity: cap,
		queue:    list.New(),
	}
}

func (s *Cache) Add(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	elem, ok := s.records[key]
	if ok {
		s.queue.MoveToFront(elem)
		record := elem.Value.(*Record)
		record.Value = value
		return
	}
	if s.queue.Len() == int(s.capacity) {
		elem = s.queue.Back()
		delete(s.records, elem.Value.(*Record).Key)
		s.queue.Remove(elem)
	}
	elem = s.queue.PushFront(&Record{Key: key, Value: value})
	s.records[key] = elem
}

func (s *Cache) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	record, ok := s.records[key]
	if !ok {
		return struct{}{}, false
	}
	return record.Value.(*Record).Value, ok
}
func (s *Cache) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	record, ok := s.records[key]
	if ok {
		s.queue.Remove(record)
		delete(s.records, key)
	}
}

func (s *Cache) Count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.records)
}

func main() {
	cache := NewCache(10)
	cache.Add("key1", "val1")
	val, ok := cache.Get("key1")
	if ok {
		fmt.Printf("%v\r\n", val)
	} else {
		fmt.Println("cache not found")
	}
}
