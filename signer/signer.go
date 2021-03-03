package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var out, in chan interface{}
	in = nil
	wg := &sync.WaitGroup{}
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
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	outerWg := &sync.WaitGroup{}
	for {
		val, ok := <-in
		if !ok {
			break
		}
		data := fmt.Sprintf("%v", val)
		fmt.Println(data, "SingleHash data", data)
		d2 := DataSignerMd5(data)
		fmt.Println(data, "SingleHash md5(data)", d2)
		outerWg.Add(1)
		go func() {
			defer outerWg.Done()
			var d1, d3 string
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				d1 = DataSignerCrc32(data)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				d3 = DataSignerCrc32(d2)
			}()
			wg.Wait()
			fmt.Println(data, "SingleHash crc32(data)", d1)
			fmt.Println(data, "SingleHash crc32(md5(data))", d3)
			result := d1 + "~" + d3
			fmt.Println(data, "SingleHash result", result)
			out <- result
		}()
	}
	outerWg.Wait()
}
func MultiHash(in, out chan interface{}) {
	outerWg := &sync.WaitGroup{}
	for {
		val, ok := <-in
		if !ok {
			break
		}
		data := fmt.Sprintf("%v", val)
		outerWg.Add(1)
		go func() {
			defer outerWg.Done()
			var result [6]string
			wg := &sync.WaitGroup{}
			for i := 0; i <= 5; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					th := strconv.Itoa(i)
					result[i] = DataSignerCrc32(th + data)
					fmt.Println(data, "MultiHash: crc32(th+step1))", i, result[i])
				}(i)
			}
			wg.Wait()
			resultString := ""
			for _, s := range result {
				resultString += s
			}
			fmt.Println(data, "MultiHash: result", resultString)
			out <- resultString
		}()
	}
	outerWg.Wait()
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
	if len(result) > 0 {

		resultString = result[0]
		for i := 1; i < len(result); i++ {

			resultString += "_" + result[i]
		}
	}
	fmt.Println("CombineResults", resultString)
	out <- resultString
}
