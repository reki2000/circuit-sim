package sim

type Static struct {
	out  int
	v    int
	name string
}

func (m *Static) Simulate(w *Wires, visited []int) []int {
	if !contains(visited, m.out) {
		w.set(m.out, m.v)
		return []int{m.out}
	}
	return []int{}
}
func (m *Static) Name() string {
	return m.name
}

type Mos struct {
	s, g, d int
	typeP   bool
	name    string
}

func (m *Mos) Simulate(w *Wires, visited []int) []int {
	s, g := w.get(m.s), w.get(m.g)
	if contains(visited, m.s) {
		if (m.typeP && g < s) || (!m.typeP && g > s) {
			if w.get(m.d) != s || !contains(visited, m.d) {
				w.set(m.d, s)
				return []int{m.d}
			}
		}
	}
	return []int{}
}

func (m *Mos) Name() string {
	return m.name
}

type Clock struct {
	vcc, gnd int
	out      int
	duration int
	name     string
	on       bool
	count    int
}

func (c *Clock) Simulate(w *Wires, visited []int) []int {
	if !contains(visited, c.out) {
		if c.on {
			w.set(c.out, w.get(c.vcc))
		} else {
			w.set(c.out, w.get(c.gnd))
		}
		c.count--
		if c.count <= 0 {
			c.count = c.duration
			c.on = !c.on
		}
		return []int{c.out}
	}
	return []int{}
}

func (c *Clock) Name() string {
	return c.name
}
