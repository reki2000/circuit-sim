package sim

import (
	"fmt"
)

type Simulation struct {
	*Circuit
	records []map[int]int
}

func NewSimulation(c *Circuit) *Simulation {
	return &Simulation{Circuit: c, records: []map[int]int{}}
}

func (s *Simulation) resetRecord() {
	s.monitee = map[int]string{}
	s.records = []map[int]int{}
}

func (s *Simulation) record() {
	r := map[int]int{}
	for k := range s.monitee {
		r[k] = s.get(k)
	}
	s.records = append(s.records, r)
}

func (s *Simulation) showRecords() (out map[string]string) {
	fmt.Printf("%20s[%03d]: %s\n", "", 0, "0123456789012345678901234567890123456789")
	out = map[string]string{}
	keys := sortIntStringMapByValue(s.monitee)
	for _, k := range keys {
		r := ""
		for t := 0; t < len(s.records); t++ {
			r += fmt.Sprintf("%s", wireValueToChar(s.records[t][k.Key]))
		}
		out[k.Value] = r
		fmt.Println(fmt.Sprintf("%20s[%03d]: ", k.Value, k.Key) + r)
	}
	return
}

func (w *Wires) formatWireFull(visited []int) string {
	result := ""
	for i := 0; i <= w.maxWireID; i++ {
		ch := " "
		if contains(visited, i) {
			ch = "*"
		}
		result += fmt.Sprintf("%3d:%d%s ->%v :%s\n", i, w.get(i), ch, w.listBonded(i), w.monitee[i])

	}
	return result + "\n"
}

func (w *Wires) formatWire(visited []int) string {
	result := ""
	for i := 0; i <= w.maxWireID; i++ {
		ch := " "
		if contains(visited, i) {
			ch = "*"
		}
		result += fmt.Sprintf("%3d:%d%s ", i, w.get(i), ch)
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

func runScenario2(c *Circuit, scenario2 map[int]string, debugName string) map[string]string {
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
	return runScenario(c, scenario, debugName)
}

func runScenario(c *Circuit, scenario map[int][]int, debugName string) map[string]string {
	play := []map[int]int{}
	for k, v := range scenario {
		for t, val := range v {
			if len(play) <= t {
				play = append(play, make(map[int]int))
			}
			play[t][k] = val
		}
	}

	s := NewSimulation(c)
	for t, set := range play {
		for k, v := range set {
			c.visit(v, k, []int{}, false)
		}
		result := c.simulateAll(debugName)
		if !result {
			fmt.Printf("#%3d: Simulation instable !!!!!!!!!!!!!!!!!!!!!!\n", t)
			s.record()
		} else {
			s.record()
		}
	}
	return s.showRecords()
}
