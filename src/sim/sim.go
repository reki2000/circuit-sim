package sim

import (
	"fmt"
	"strings"
)

const numWires = 100

var wire [numWires]int
var maxWireID = -1

func w() int {
	if maxWireID >= len(wire)-1 {
		panic("too many wires")
	}
	maxWireID++
	return maxWireID
}

var net struct {
	from [numWires][]int
	to   [numWires][]int
}

func bond(w1, w2 int) {
	//fmt.Printf(" bonding %d to %d\n", w1, w2)
	if net.from[w1] != nil {
		net.from[w1] = append(net.from[w1], w2)
	} else {
		net.from[w1] = []int{w2}
	}
}

func listBonded(w int) []int {
	from, to := net.from[w], net.to[w]
	if from == nil {
		return to
	} else if to == nil {
		return from
	}
	return append(net.from[w], net.to[w]...)
}

type Simulatable interface {
	Simulate([]int) []int
	Name() string
}

func visit(fillValue int, start int, visited []int, debug bool) []int {
	waiting := []int{start}
	str := fmt.Sprintf(" filling wires with %d: ", fillValue)
	for w := 0; len(waiting) > 0; {
		w, waiting = pop(waiting)
		wire[w] = fillValue
		str += fmt.Sprintf("%3d ", w)
		visited = push(visited, w)
		for _, nextWire := range listBonded(w) {
			waiting = push(waiting, nextWire)
		}
	}
	if debug {
		fmt.Println(str)
	}
	return visited
}

func simulateAll(debugName string) bool {
	visited := []int{}
	updated := true
	for loop := 0; loop < 5 && updated; loop++ {
		if debugName != "" {
			fmt.Printf("Simulation #%d start ..............................\n%v\n", loop, net)
		}
		updated = false
		for _, m := range modules {
			outWires := m.Simulate(visited)
			debug := false
			if debugName != "" && strings.HasPrefix(m.Name(), debugName) {
				fmt.Printf("dev[%v] out:%v wires:\n%v\n", m, outWires, formatWire(visited))
				debug = true
			}
			for _, out := range outWires {
				updated = true
				fillValue := wire[out]
				visited = visit(fillValue, out, visited, debug)
			}
		}
	}
	if updated {
		return false
	}
	return true
}

var modules []Simulatable

func addModule(m Simulatable) {
	modules = append(modules, m)
}

func formatWire(visited []int) string {
	result := ""
	for k, v := range monitee {
		result += fmt.Sprintf("%20s[%3d]:%d\n", v, k, wire[k])
	}
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

func Test() {
}
