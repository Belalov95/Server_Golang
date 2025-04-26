package main

import (
	"fmt"
	"sort"
	"testing"
)

// Задача из гита
func Test(t *testing.T) {
	slice := []int{1, 2, 3}
	slice = append(slice, 4)
	fmt.Println(slice[len(slice)-1]) // 4
}

// задача №1
func main() {
	v := []int{3, 4, 1, 2, 5}
	ap(v)
	sr(v)
	fmt.Println(v) //1, 2, 3, 4, 5
	// тут 10 не добавляется, т.к. у нас в функции "ap" емкость равна 5, из-за нехватки места
	// создается новый массив, а в функцию "main" мы не передаем этот новый массив и,соответсвенно,
	// у нас наш срез v не меняется
}
func ap(arr []int) {
	arr = append(arr, 10)
}
func sr(arr []int) {
	sort.Ints(arr)
}

// задача №2
func Main() {
	var foo []int
	var bar []int

	foo = append(foo, 1)
	foo = append(foo, 2)
	foo = append(foo, 3)
	bar = append(foo, 4)
	foo = append(foo, 5)

	fmt.Println(foo, bar) //1, 2, 3, 5 //1, 2, 3, 4

}

// задача №3
func mAin() {
	c := []string{"A", "B", "D", "E"}
	b := c[1:2]
	b = append(b, "TT")
	fmt.Println(c) // A, B, D, E
	fmt.Println(b) // B, TT
}

// задача №4
func maun() {
	var m map[string]int
	for _, word := range []string{"hello", "world", "from", "the",
		"best", "language", "in", "the", "world"} {
		m[word]++
	}
	for k, v := range m {
		fmt.Println(k, v)
		// hello: 1
		// world: 2
		// from: 1
		// the: 2
		// best: 1
		// language: 1
		// in: 1
	}
}

// другое написание кода для 4 задачи
func example() {
	m := make(map[string]int)
	words := []string{"hello", "world", "from", "the", "best", "language", "in", "the", "world"}
	for _, word := range words {
		m[word]++
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}
