package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Notify struct {
}

type Keylistener interface {
	Stop()
	Keys() <-chan int
}

type ReadKey struct {
	ch chan int
}

func (r *ReadKey) Stop() {

}

func (r *ReadKey) Keys() <-chan int {
	return r.ch
}

func (n *Notify) GetNotification() Keylistener {
	read := ReadKey{ch: make(chan int)}
	//Добавить горутины которая каждую секунду будет писать в канал read.ch рандомное значение от нуля до ста.
	//надо в main записать 5 значений из этого канала и сделать fmt.Print
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				value := rand.IntN(101)
				read.ch <- value
			}
		}
	}()

	return &read
}

func main() {

	notify := Notify{}
	key := notify.GetNotification()
	values := make([]int, 0, 5)
	for n := range key.Keys() {
		values = append(values, n)
		if len(values) == 5 {
			break
		}
	}

	fmt.Println(values)
}
