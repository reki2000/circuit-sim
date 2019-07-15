package sim

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
