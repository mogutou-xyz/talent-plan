package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

// URLCountMap is the map function in the first round
func URLCountMap(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if len(l) == 0 {
			continue
		}
		kvs = append(kvs, KeyValue{Key: l})
	}
	return kvs
}

// URLCountReduce is the reduce function in the first round
func URLCountReduce(key string, values []string) string {
	return fmt.Sprintf("%s %s\n", key, strconv.Itoa(len(values)))
}

// URLTop10Map is the map function in the second round
// 优化只取 前十名
func URLTop10Map(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	us, cs := urlTop10(lines)
	kvs := make([]KeyValue, 0, len(us))
	for k, v := range us {
		kvs = append(kvs, KeyValue{Value: fmt.Sprintf("%s %d\n", v, cs[k])})
	}

	return kvs
}

// URLTop10Reduce is the reduce function in the second round
func URLTop10Reduce(key string, values []string) string {
	us, cs := urlTop10(values)
	buf := new(bytes.Buffer)
	for i := range us {
		fmt.Fprintf(buf, "%s: %d\n", us[i], cs[i])
	}
	return buf.String()
}

func urlTop10(values []string) ([]string, []int) {
	cnts := make(map[string]int, len(values))
	for i := 0; i < len(values); i++ {
		v := strings.TrimSpace(values[i])
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		cnts[tmp[0]] = n
	}

	return TopN(cnts, 10)
}
