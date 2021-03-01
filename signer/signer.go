package main

import (
	"fmt"
	"sync"
	"time"
)

func ExecutePipeline(jobs ...job) {
	var out, in chan interface{}
	in = nil
	wg := &sync.WaitGroup{} // wait_2.go инициализируем группу
	for i := 0; i < len(jobs); i++ {
		out = make(chan interface{})
		wg.Add(1)
		i := i
		go func(in, out chan interface{}) {
			defer wg.Done()
			jobs[i](in, out)
		}(in, out)
		in = out
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}

func transmitter(in, out chan int) {
	for {
		val, ok := <-in
		val++
		out <- val
		if !ok {
			break
		}
	}
}

func main() {
	in := make(chan int)
	firstIn := in
	var out chan int
	for i := 0; i < 5; i++ {
		out = make(chan int)
		go transmitter(in, out)
		in = out
	}
	upper := 300
	for i := 100; i <= upper-1; i += 100 {
		firstIn <- i
	}
	fmt.Scanln()
	for i := range out {
		fmt.Println("iout", i)
	}
	firstIn <- 1
	//out = make(chan int)
	//go transmitter(in, out)
	//in = out
	//out = make(chan int)
	//go transmitter(in, out)
	fmt.Print("out", <-out)

}

// сюда писать код
