package sim

import (
	"fmt"
	"sort"
)

var monitee = map[int]string{}
var records = []map[int]int{}

func monitor(targets map[int]string) {
	for k, v := range targets {
		monitee[k] = v
	}
}

func record() {
	r := map[int]int{}
	for k := range monitee {
		r[k] = wire[k]
	}
	records = append(records, r)
}

func showRecords() (out map[string]string) {
	fmt.Printf("%20s[%03d]: %s\n", "", 0, "0123456789012345678901234567890123456789")
	out = map[string]string{}
	keys := []int{}
	for k := range monitee {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		r := ""
		for t := 0; t < len(records); t++ {
			r += fmt.Sprintf("%s", wireValueToChar(records[t][k]))
		}
		out[monitee[k]] = r
		fmt.Println(fmt.Sprintf("%20s[%03d]: ", monitee[k], k) + r)
	}
	return
}
