package main

import (
	"fmt"
	"sort"
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
			close(out)
		}(in, out)
		in = out
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}

var singleHashCalled int = 0

func SingleHash(in, out chan interface{}) {
	for {
		val, ok := <-in
		if !ok {
			break
		}
		fmt.Println(val, "SingleHash val")
		data := fmt.Sprintf("%v", val)
		fmt.Println(data, "SingleHash data", data)
		d1 := DataSignerCrc32(data)
		fmt.Println(data, "SingleHash crc32(data)", d1)
		d2 := DataSignerMd5(data)
		fmt.Println(data, "SingleHash md5(data)", d2)
		d3 := DataSignerCrc32(d2)
		fmt.Println(data, "SingleHash crc32(md5(data))", d3)
		result := d1 + "~" + d3
		fmt.Println(data, "SingleHash result", result)
		out <- result
		singleHashCalled++
	}
}
func MultiHash(in, out chan interface{}) {
	for {
		val, ok := <-in
		if !ok {
			break
		}
		fmt.Println(val, "MultiHash val")
		data := fmt.Sprintf("%v", val)
		result := ""
		for i := 0; i <= 5; i++ {
			i := i
			tempRes := DataSignerCrc32(string(i) + data)
			fmt.Println(data, "MultiHash: crc32(th+step1))", i, tempRes)
			result += tempRes
		}
		fmt.Println(data, "MultiHash: result", result)
		out <- result
	}
}
func CombineResults(in, out chan interface{}) {
	var result []string
	for {
		val, ok := <-in
		if !ok {
			break
		}
		result = append(result, fmt.Sprintf("%v", val))

	}
	sort.Sort(sort.StringSlice(result))
	resultString := ""
	for _, s := range result {

		resultString += "_" + s
	}
	fmt.Println("CombineResults", result)
	out <- result
}
