package sim

import (
	"fmt"
	"strings"
)

type Simulatable interface {
	Simulate([]int) []int
	Name() string
}

var devices []Simulatable

func addDevice(m Simulatable) {
	devices = append(devices, m)
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

const simulationLoopCount = 10

func simulateAll(debugName string) bool {
	visited := []int{}
	updated := true
	for loop := 0; loop < simulationLoopCount && updated; loop++ {
		if debugName != "" {
			fmt.Printf("Simulation #%d start ..............................\n%v\n", loop, formatWireFull(visited))
		}
		updated = false
		for _, m := range devices {
			outWires := m.Simulate(visited)
			debug := false
			if (debugName != "" && strings.HasPrefix(m.Name(), debugName)) || debugName == "*" {
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

func Test() {
}
