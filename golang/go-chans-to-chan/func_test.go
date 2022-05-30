package func_test

import (
	"sort"
	"testing"
    "time"
)

func TestCombineOne(t *testing.T) {
	var (
		requiredNum = 30
		result      = []int{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		}
		chanOne = make(chan int)
		chanTwo = make(chan int)
		chanThr = make(chan int)
	)

	generate := func() {
		go sendData(0, 10, time.Second*1, 1, chanOne)
		go sendData(10, 20, time.Second*2, 2, chanTwo)
		go sendData(20, 30, time.Second*3, 3, chanThr)
	}

	generate()
	res := consume(combine(chanOne, chanTwo, chanThr))
	if len(res) != 30 {
		t.Fatalf("result size is not eq %d", requiredNum)
	}
	sort.Ints(res)
	for i := range res {
		if res[i] != result[i] {
			t.Fatal("slices are not equal")
		}
	}
}

func BenchmarkCombineOne(b *testing.B) {
	var (
		chanOne = make(chan int)
		chanTwo = make(chan int)
		chanThr = make(chan int)
	)

	generate := func() {
		go sendData(0, 10, time.Second*1, 1, chanOne)
		go sendData(10, 20, time.Second*2, 2, chanTwo)
		go sendData(20, 30, time.Second*3, 3, chanThr)
	}
	generate()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			combine(chanOne, chanTwo, chanThr)
		}
	})
}

func BenchmarkCombineTwo(b *testing.B) {
	var (
		chanOne = make(chan int)
		chanTwo = make(chan int)
		chanThr = make(chan int)
	)

	generate := func() {
		go sendData(0, 10, time.Second*1, 1, chanOne)
		go sendData(10, 20, time.Second*2, 2, chanTwo)
		go sendData(20, 30, time.Second*3, 3, chanThr)
	}
	generate()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			combineTwo(chanOne, chanTwo, chanThr)
		}
	})
}
