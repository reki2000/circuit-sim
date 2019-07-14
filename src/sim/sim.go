package sim

import (
	"fmt"
	"sort"
	"strings"
)

var wire [100]int

var net struct {
	from [100][]int
	to   [100][]int
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
	Name() string
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
func (m *Static) Name() string {
	return m.name
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

func (m *Mos) Name() string {
	return m.name
}

var modules []Simulatable

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

func simulateAll(debugName string) {
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
		fmt.Println("Simulation instable !!!!!!!!!!!!!!!!!!!!!!")
	}
}

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

func showRecords() (out map[string]string) {
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
	for _, set := range play {
		for k, v := range set {
			visit(v, k, []int{}, false)
		}
		simulateAll(debugName)
		record()
	}
	return showRecords()
}

var maxWireID = -1

func w() int {
	if maxWireID >= len(wire)-1 {
		panic("too many wires")
	}
	maxWireID++
	return maxWireID
}

func setup() {
	gnd, vdd := w(), w()
	addModule(&Static{gnd, 0, "Gnd"})
	addModule(&Static{vdd, 1, "Vdd"})
}

func Test() {
	setup()

	clk, d, q := buildGatedDLatch("d0")

	// wireID: { value series for each t ...}
	scenario := map[int][]int{
		clk: {1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0},
		d:   {1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	runScenario(scenario, "")
}
