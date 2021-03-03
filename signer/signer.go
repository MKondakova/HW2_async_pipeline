package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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
		md5Hash := DataSignerMd5(data)
		outerWg.Add(1)
		go func() {
			defer outerWg.Done()
			var crc32FromData, crc32FromMd5 string
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				crc32FromData = DataSignerCrc32(data)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				crc32FromMd5 = DataSignerCrc32(md5Hash)
			}()
			wg.Wait()
			result := crc32FromData + "~" + crc32FromMd5
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
				}(i)
			}
			wg.Wait()
			out <- strings.Join(result[:], "")
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
	out <- strings.Join(result, "_")
}
