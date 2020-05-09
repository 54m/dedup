package main

import (
	"fmt"

	"github.com/54mch4n/dedup"
)

func main() {
	ded := dedup.NewDeduplication()
	sliceStr := []string{"Go", "V", "Java", "Python", "Go", "Ruby", "Go", "V"}
	fmt.Println(ded.Do(sliceStr).Str())

	sliceInt := []int{1, 1, 2, 2, 3, 3}
	fmt.Println(ded.Do(sliceInt).Int())

	sliceFloat64 := []float64{0.1, 0.1, 0.2, 0.2, 0.3, 0.3}
	fmt.Println(ded.Do(sliceFloat64).Float64())
}
