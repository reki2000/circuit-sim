package sim

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

type Mos struct {
	s, g, d int
	typeP   bool
	name    string
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

type Clock struct {
	vcc, gnd int
	out      int
	duration int
	name     string
	on       bool
	count    int
}

func (c *Clock) Simulate(visited []int) []int {
	if !contains(visited, c.out) {
		if c.on {
			wire[c.out] = wire[c.vcc]
		} else {
			wire[c.out] = wire[c.gnd]
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
