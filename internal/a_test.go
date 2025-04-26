package main

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

/*
	func TestXxx(t *testing.T) {
		a := make([]int, 0, 4)
		a = append(a, 1, 2, 3)
		fmt.Println(a, cap(a), len(a)) //1, 2, 3
		foo(a)
		fmt.Println(a[:cap(a)], cap(a), len(a)) // 1, 2, 3
	}

	func foo(b []int) {
		b = append(b, 4)
		fmt.Println(b, cap(b), len(b)) //1, 2, 3, 4
	}
*/
func TestXxx(t *testing.T) {
	k := Store{m: map[string]int{}}
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(k.Get(strconv.Itoa(i)))
		}()
		go func() {
			defer wg.Done()
			k.Set(strconv.Itoa(i), i)
		}()
	}
	wg.Wait()
}

type Store struct {
	m map[string]int
}

func (s *Store) Get(key string) (int, bool) {
	val, ok := s.m[key]
	return val, ok
}
func (s *Store) Set(key string, value int) {
	s.m[key] = value
}
