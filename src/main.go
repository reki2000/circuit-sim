package main

import (
	"fmt"
	"sort"
	"strconv"
)

var wire [100]int

var net struct {
	from [100][]int
	to   [100][]int
}

func bond(w1, w2 int) {
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

type Simulatable interface {
	Simulate([]int) []int
}

type Mos struct {
	s, g, d int
	typeP   bool
	name    string
}

type Static struct {
	out  int
	v    int
	name string
}

func (m *Static) Simulate(visited []int) []int {
	if !contains(visited, m.out) {
		wire[m.out] = m.v
		return []int{m.out}
	}
	return []int{}
}

func (m *Mos) Simulate(visited []int) []int {
	s, g := wire[m.s], wire[m.g]
	if contains(visited, m.s) {
		if (m.typeP && g < s) || (!m.typeP && g > s) {
			if wire[m.d] != wire[m.s] || !contains(visited, m.d) {
				wire[m.d] = wire[m.s]
				return []int{m.d}
			}
		}
	}
	return []int{}
}

var modules []Simulatable

func visit(fillValue int, start int, visited []int, debug bool) []int {
	waiting := []int{start}
	for w := 0; len(waiting) > 0; {
		w, waiting = pop(waiting)
		wire[w] = fillValue
		visited = push(visited, w)
		if debug {
			fmt.Printf(" wire[%d] = %d\n", w, fillValue)
		}
		for _, nextWire := range listBonded(w) {
			waiting = push(waiting, nextWire)
		}
	}
	return visited
}

func simulateAll(debug bool) {
	visited := []int{}
	updated := true
	for loop := 0; loop < 5 && updated; loop++ {
		updated = false
		for _, m := range modules {
			outWires := m.Simulate(visited)
			if debug {
				fmt.Printf("dev[%v] %v wires:%v\n", m, outWires, wire[0:5])
			}
			for _, out := range outWires {
				updated = true
				fillValue := wire[out]
				visited = visit(fillValue, out, visited, debug)
			}
		}
	}
}

func addModule(m Simulatable) {
	modules = append(modules, m)
}

func showWire(max int) {
	result := ""
	for k, v := range monitee {
		result += fmt.Sprintf("%20s[%d]:%d\n", v, k, wire[k])
	}
	fmt.Println(result)
}

var monitee = map[int]string{}
var records = []map[int]int{}

func monitor(targets map[int]string) {
	for k, v := range targets {
		monitee[k] = v
	}
}

func record() {
	r := map[int]int{}
	for k, _ := range monitee {
		r[k] = wire[k]
	}
	records = append(records, r)
}

func wireValueToChar(value int) string {
	var ch string
	switch value {
	case 0:
		ch = "_"
	case 1:
		ch = "H"
	default:
		ch = "?"
	}
	return ch
}

func showRecords() {
	keys := []int{}
	for k, _ := range monitee {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		r := fmt.Sprintf("%20s[%03d]: ", monitee[k], k)
		for t := 0; t < len(records); t++ {
			r += fmt.Sprintf("%s", wireValueToChar(records[t][k]))
		}
		fmt.Println(r)
	}
}

func runScenario(scenario map[int][]int, debug bool) {
	play := []map[int]int{}
	for k, v := range scenario {
		for t, val := range v {
			if len(play) <= t {
				play = append(play, make(map[int]int))
			}
			play[t][k] = val
		}
	}
	for _, set := range play {
		for k, v := range set {
			visit(v, k, []int{}, false)
		}
		simulateAll(debug)
		record()
	}
	showRecords()
}

var maxWireId = -1

func w() int {
	if maxWireId >= len(wire)-1 {
		panic("too many wires")
	}
	maxWireId++
	return maxWireId
}

func buildNandGate(name string) (in1, in2, out int) {
	in1, in2, out = w(), w(), w()
	out2 := w()
	addModule(&Mos{1, in1, out, true, name + ".p0"})
	addModule(&Mos{1, in2, out, true, name + ".p1"})
	addModule(&Mos{0, in2, out2, false, name + ".n1"})
	addModule(&Mos{out2, in1, out, false, name + ".n0"})
	return
}

func buildNotGate(name string) (in, out int) {
	in, out = w(), w()
	addModule(&Mos{1, in, out, true, name + ".p"})
	addModule(&Mos{0, in, out, false, name + ".n"})
	return
}

func buildDFlipFlop(name string) (clk, d, q int) {
	clk, d, q = w(), w(), w()
	in1, in2, out := [8]int{}, [8]int{}, [8]int{}
	for i := 0; i < 8; i++ {
		myName := name + ".nand" + strconv.Itoa(i)
		in1[i], in2[i], out[i] = buildNandGate(myName)
		monitor(map[int]string{
			in1[i]: myName + ".in1",
			in2[i]: myName + ".in2",
			out[i]: myName + ".out",
		})
	}
	noti, noto := [3]int{}, [3]int{}
	for i := 0; i < 3; i++ {
		noti[i], noto[i] = buildNotGate(name + ".not" + string(i))
	}

	// clk ot not0
	bond(clk, noti[0])

	// not0 to nand0, nand1, not1
	bond(noto[0], in2[0])
	bond(noto[0], in2[1])
	bond(noto[0], noti[1])

	// d to nand0, not2
	bond(d, in1[0])
	bond(d, noti[2])

	// not2 to nand1
	bond(noto[2], in1[1])

	// first RS latch
	bond(out[0], in1[2])
	bond(out[1], in2[3])
	bond(out[2], in1[3])
	bond(out[3], in2[2])

	// first RS latch to second
	bond(out[2], in1[4])
	bond(out[3], in1[5])

	// not1 to nadn4, nand5
	bond(noto[1], in2[4])
	bond(noto[1], in2[5])

	// final RS latch
	bond(out[4], in1[6])
	bond(out[5], in2[7])
	bond(out[6], in2[7])
	bond(out[7], in2[6])

	// output
	bond(out[6], q)
	return
}

func testDFlipFlop() {
	clk, d, q := buildDFlipFlop("d0")

	// wireID: { value series for each t ...}
	scenario := map[int][]int{
		clk: {1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0},
		d:   {1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	runScenario(scenario, false)
}

func testRSLatch() {
	r, s, q := w(), w(), w()
	i11, i12, o1 := buildNandGate("nand1")
	i21, i22, o2 := buildNandGate("nand2")
	bond(s, i11)
	bond(o2, i12)
	bond(r, i21)
	bond(o1, i22)
	bond(o1, q)
	monitor(map[int]string{r: "R", s: "S", q: "Q", i11: "i11", i12: "i12", o1: "o1", i21: "i21", i22: "i22", o2: "o2"})

	scenario := map[int][]int{
		r: {1, 1, 1, 1, 0, 1, 1, 1},
		s: {1, 0, 1, 1, 1, 1, 0, 1},
	}
	runScenario(scenario, false)
}

func main() {
	gnd, vdd := w(), w()
	addModule(&Static{gnd, 0, "Gnd"})
	addModule(&Static{vdd, 1, "Vdd"})

	//testRSLatch()
	testDFlipFlop()
}
