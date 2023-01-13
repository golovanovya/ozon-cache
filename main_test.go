package main

import (
	"sync"
	"testing"
)

func TestCache_Add_existElementWithFullQueueSync_moveToFront(t *testing.T) {
	cache := NewCache(3)
	cache.Add("key1", "val1")
	cache.Add("key2", "val2")
	cache.Add("key3", "val3")
	elem := cache.queue.Back()
	if elem.Value.(*Record).Value.(string) != "val1" {
		t.Fatal("неправильное значение")
	}
	cache.Add("key1", "val3")
	elem = cache.queue.Front()
	if elem.Value.(*Record).Value.(string) != "val3" {
		t.Fatal("неправильное значение")
	}
	if cache.Count() > 3 {
		t.Fatal("длина кэша больше вместимости")
	}
}

func TestCaceh_ParallelAddElement(t *testing.T) {
	cache := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		cache.Add("key1", "value1")
	}()
	go func() {
		defer wg.Done()
		cache.Add("key2", "value2")
	}()
	go func() {
		defer wg.Done()
		cache.Add("key3", "value3")
	}()
	wg.Wait()
	val, ok := cache.Get("key1")
	if !ok {
		t.Fatal("значение не найдено")
	}
	if val.(string) != "value1" {
		t.Fatal("значение не корректно")
	}
	val, ok = cache.Get("key2")
	if !ok {
		t.Fatal("значение не найдено")
	}
	if val.(string) != "value2" {
		t.Fatal("значение не корректно")
	}
	val, ok = cache.Get("key3")
	if !ok {
		t.Fatal("значение не найдено")
	}
	if val.(string) != "value3" {
		t.Fatal("значение не корректно")
	}
}

func TestCache_AddElement_CorrectAddingElement(t *testing.T) {
	cache := NewCache(10)
	cache.Add("key1", "val1")
	val, ok := cache.Get("key1")
	if !ok {
		t.Fatal("элемент не найден в кэше")
	}
	if val.(string) != "val1" {
		t.Fatal("значение элемента неверно")
	}
}

func TestCache_Get_NotExists_ReturnFalse(t *testing.T) {
	cache := NewCache(10)
	cache.Add("key1", "val1")
	_, ok := cache.Get("key2")
	if ok != false {
		t.Fatal("должен не быть найден")
	}
}

func TestCache_Delete_DeleteExistsElement(t *testing.T) {
	cache := NewCache(10)
	cache.Add("key1", "val1")
	if cache.Count() != 1 {
		t.Fatal("элемент не добавлен")
	}
	cache.Delete("key1")
	if cache.Count() > 0 {
		t.Fatal("элемент не удалён")
	}
}

func TestCache_Count_OverflowCapacity_ReturnCorrectCount(t *testing.T) {
	cache := NewCache(2)
	cache.Add("key1", "val1")
	cache.Add("key2", "val2")
	cache.Add("key3", "val3")
	if cache.Count() > 2 {
		t.Fatal("количество элементов больше вместимости")
	}
}
