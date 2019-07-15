package sim

import (
	"fmt"
	"strings"
)

type Circuit struct {
	*Wires
	devices  []Simulatable
	gnd, vdd int
}

type Simulatable interface {
	Simulate(*Wires, []int) []int
	Name() string
}

func NewCircuit() *Circuit {
	c := &Circuit{Wires: NewWire(), devices: make([]Simulatable, 0)}
	c.gnd, c.vdd = c.w(), c.w()

	c.addDevice(&Static{c.gnd, 0, "Gnd"})
	c.addDevice(&Static{c.vdd, 1, "Vdd"})

	return c
}

func (c *Circuit) addDevice(m Simulatable) {
	c.devices = append(c.devices, m)
}

func (c *Circuit) visit(fillValue int, start int, visited []int, debug bool) []int {
	waiting := []int{start}
	str := fmt.Sprintf(" filling wires with %d: ", fillValue)
	for w := 0; len(waiting) > 0; {
		w, waiting = pop(waiting)
		c.set(w, fillValue)
		str += fmt.Sprintf("%3d ", w)
		visited = push(visited, w)
		for _, nextWire := range c.listBonded(w) {
			waiting = push(waiting, nextWire)
		}
	}
	if debug {
		fmt.Println(str)
	}
	return visited
}

const simulationLoopCount = 5

func (c *Circuit) simulateAll(debugName string) bool {
	visited := []int{}
	updated := true
	for loop := 0; loop < simulationLoopCount && updated; loop++ {
		if debugName != "" {
			fmt.Printf("Simulation #%d start ........................\n%v\n", loop, c.formatWireFull(visited))
		}
		updated = false
		for _, m := range c.devices {
			outWires := m.Simulate(c.Wires, visited)
			debug := false
			if (debugName != "" && strings.HasPrefix(m.Name(), debugName)) || debugName == "*" {
				fmt.Printf("dev[%v] out:%v wires:\n%v\n", m, outWires, c.formatWire(visited))
				debug = true
			}
			for _, out := range outWires {
				updated = true
				fillValue := c.get(out)
				visited = c.visit(fillValue, out, visited, debug)
			}
		}
	}
	if updated {
		return false
	}
	return true
}
