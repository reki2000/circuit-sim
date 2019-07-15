package sim

import "fmt"

type Wires struct {
	wire      []int
	maxWireID int
	from      [][]int
	monitee   map[int]string
}

func NewWire() *Wires {
	numWires := 1000
	return &Wires{make([]int, numWires), -1, make([][]int, numWires), map[int]string{}}
}

func (w *Wires) w() int {
	if w.maxWireID >= len(w.wire)-1 {
		panic("too many wires")
	}
	w.maxWireID++
	return w.maxWireID
}

func (w *Wires) set(id, val int) {
	w.wire[id] = val
}

func (w *Wires) get(id int) int {
	return w.wire[id]
}

func (w *Wires) bond(w1, w2 int) {
	//fmt.Printf(" bonding %d to %d\n", w1, w2)
	if w.from[w1] != nil {
		w.from[w1] = append(w.from[w1], w2)
	} else {
		w.from[w1] = []int{w2}
	}
}

func (w *Wires) bondWires(w1, w2 []int) {
	if len(w1) != len(w2) {
		panic(fmt.Sprintf("bonding length mismatched wires:%v to %v\n", w1, w2))
	}
	for i, v := range w1 {
		w.bond(v, w2[i])
	}
}

func (w *Wires) listBonded(id int) []int {
	return w.from[id]
}

func (w *Wires) monitor(id int, name string) {
	w.monitee[id] = name
}

func (w *Wires) monitor2(targets map[int]string) {
	for k, v := range targets {
		w.monitor(k, v)
	}
}
