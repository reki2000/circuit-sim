package sim

import (
	"fmt"
)

var monitee map[int]string
var records []map[int]int

func monitor(targets map[int]string) {
	for k, v := range targets {
		monitee[k] = v
	}
}
func resetRecord() {
	monitee = map[int]string{}
	records = []map[int]int{}
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
	keys := sortIntStringMapByValue(monitee)
	for _, k := range keys {
		r := ""
		for t := 0; t < len(records); t++ {
			r += fmt.Sprintf("%s", wireValueToChar(records[t][k.Key]))
		}
		out[k.Value] = r
		fmt.Println(fmt.Sprintf("%20s[%03d]: ", k.Value, k.Key) + r)
	}
	return
}

func formatWireFull(visited []int) string {
	result := ""
	for i := 0; i <= maxWireID; i++ {
		ch := " "
		if contains(visited, i) {
			ch = "*"
		}
		result += fmt.Sprintf("%3d:%d%s ->%v :%s\n", i, wire[i], ch, listBonded(i), monitee[i])

	}
	return result + "\n"
}

func formatWire(visited []int) string {
	result := ""
	for i := 0; i <= maxWireID; i++ {
		ch := " "
		if contains(visited, i) {
			ch = "*"
		}
		result += fmt.Sprintf("%3d:%d%s ", i, wire[i], ch)
	}
	return result + "\n"
}

func wireValueToChar(value int) string {
	var ch string
	switch value {
	case 0:
		ch = "."
	case 1:
		ch = "H"
	default:
		ch = "?"
	}
	return ch
}

func runScenario2(scenario2 map[int]string, debugName string) map[string]string {
	scenario := map[int][]int{}
	for k, v := range scenario2 {
		newValues := make([]int, len(v))
		for i, ch := range v {
			val := 0
			if ch == 'H' {
				val = 1
			}
			newValues[i] = val
		}
		scenario[k] = newValues
	}
	return runScenario(scenario, debugName)
}

func runScenario(scenario map[int][]int, debugName string) map[string]string {
	play := []map[int]int{}
	for k, v := range scenario {
		for t, val := range v {
			if len(play) <= t {
				play = append(play, make(map[int]int))
			}
			play[t][k] = val
		}
	}

	for t, set := range play {
		for k, v := range set {
			visit(v, k, []int{}, false)
		}
		result := simulateAll(debugName)
		if !result {
			fmt.Printf("#%3d: Simulation instable !!!!!!!!!!!!!!!!!!!!!!\n", t)
			record()
		} else {
			record()
		}
	}
	return showRecords()
}
