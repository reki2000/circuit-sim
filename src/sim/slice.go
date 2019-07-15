package sim

import "sort"

func pop(in []int) (int, []int) {
	if len(in) == 0 {
		panic("pop from empty")
	}
	if len(in) == 1 {
		return in[0], []int{}
	}
	return in[0], in[1:]
}

func contains(in []int, v int) bool {
	for _, val := range in {
		if val == v {
			return true
		}
	}
	return false
}

func push(in []int, v int) []int {
	if !contains(in, v) {
		return append(in, v)
	}
	return in
}

type IntStringPair struct {
	Key   int
	Value string
}

func sortIntStringMapByValue(m map[int]string) []IntStringPair {
	var ss []IntStringPair
	for k, v := range m {
		ss = append(ss, IntStringPair{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})
	return ss
}
