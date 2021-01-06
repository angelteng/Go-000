package main

import (
	"fmt"
	"sync"
	"time"
)

type Windows struct {
	Data []sync.Map
	Head int
	Tail int
	Cap  int
	Time int64
}

func NewWindows(cap int) *Windows {
	cap++
	return &Windows{
		Data: make([]sync.Map, cap),
		Cap:  cap,
		Head: 0,
		Tail: 0,
		Time: 0,
	}
}
func (w *Windows) Put(key string, val int) {
	nowTime := time.Now().Unix()
	if w.Time != 0 && nowTime != w.Time {
		w.Tail = (w.Tail + 1) % w.Cap
		w.Time = nowTime
		if (w.Tail+1)%w.Cap == w.Head {
			w.Pop()
		}
	} else if w.Time == 0 {
		w.Time = nowTime
	}
	v, ok := w.Data[w.Tail].LoadOrStore(key, val)
	if ok {
		w.Data[w.Tail].Store(key, v.(int)+val)
	}
}

func (w *Windows) Pop() {
	if w.Tail == w.Head {
		return
	}
	// data := w.Data[w.Head]
	w.Data[w.Head] = sync.Map{}
	w.Head = (w.Head + 1) % w.Cap
	return
}

func (w *Windows) GetTotal() map[string]int {
	result := map[string]int{}
	for i := 0; i < w.Cap; i++ {
		index := (w.Head + i) % w.Cap
		w.Data[index].Range(func(key, val interface{}) bool {
			k := key.(string)
			v := val.(int)
			if _, ok := result[k]; !ok {
				result[k] = 0
			}
			result[k] += v
			return true
		})
	}
	return result
}

func main() {
	win := NewWindows(5)
	go func() {
		for {
			win.Put("test", 1)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	for {
		win.Put("test", 1)
		fmt.Printf("%+v \n", win.Data)
		fmt.Printf("total: %+v \n", win.GetTotal())
		time.Sleep(500 * time.Millisecond)
	}
}
