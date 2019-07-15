package sim

import "testing"

func TestRSLatch(t *testing.T) {
	ci := NewCircuit()

	s, r, q, _ := ci.buildRSLatch("test")
	ci.monitor(q, "Q")

	scenario2 := map[int]string{
		r: "HHHHH.HHHHH",
		s: "H.HHHHHHH.H",
	}
	result := runScenario2(ci, scenario2, "")

	if result["Q"] != "HHHHH....HH" {
		t.Error(result)
	}
}

func TestGatedDLatch(t *testing.T) {
	ci := NewCircuit()

	clk, d, q := ci.buildGatedDLatch("d0")
	ci.monitor(q, "Q")

	scenario2 := map[int]string{
		clk: "H...H.H....HHH.H.HHH.HH...",
		d:   "H.......HH.H......H..H....",
	}

	result := runScenario2(ci, scenario2, "")

	if result["Q"] != "HHHH.......H......H..H...." {
		t.Error()
	}
}

func TestDFlipFlop(t *testing.T) {
	ci := NewCircuit()

	clk, d, q := ci.buildDFlipFlop("d")
	ci.monitor(q, "Q")

	scenario2 := map[int]string{
		clk: ".H..H.H....HHH.H.HHH.HH...",
		d:   "HH......HHHH......H..H....",
	}

	result := runScenario2(ci, scenario2, "")

	if result["Q"] != "HHHH.......HHHH..........." {
		t.Error()
	}
}

func TestClock(t *testing.T) {
	ci := NewCircuit()

	clk := ci.buildClock("CLK", 10)
	ci.monitor(clk, "CLK")

	dummy := ci.w()
	scenario := map[int]string{
		dummy: "............................................................",
	}
	result := runScenario2(ci, scenario, "")

	if result["CLK"] != "..........HHHHHHHHHH..........HHHHHHHHHH..........HHHHHHHHHH" {
		t.Error(result)
	}
}

func TestHalfAdder(t *testing.T) {
	ci := NewCircuit()

	a, b, c, s := ci.buildHalfAdder("test")
	ci.monitor(c, "C")
	ci.monitor(s, "S")

	scenario := map[int]string{
		a: "..HH",
		b: ".H.H",
	}
	result := runScenario2(ci, scenario, "")

	if result["C"] != "...H" || result["S"] != ".HH." {
		t.Error(result)
	}
}

func TestFullAdder(t *testing.T) {
	ci := NewCircuit()

	a, b, x, c, s := ci.buildFullAdder("test")
	ci.monitor(c, "C")
	ci.monitor(s, "S")

	scenario := map[int]string{
		a: "..HH..HH",
		b: ".H.H.H.H",
		x: "....HHHH",
	}
	result := runScenario2(ci, scenario, "")

	if result["C"] != "...H.HHH" || result["S"] != ".HH.H..H" {
		t.Error(result)
	}
}

func TestNbitAdder(t *testing.T) {
	ci := NewCircuit()

	a, b, c, s := ci.buildNbitAdder("test", 2)
	ci.monitor(c, "C")
	ci.monitor(s[0], "S0")
	ci.monitor(s[1], "S1")

	scenario := map[int]string{
		a[0]: ".H.H.H.H.H.H.H.H",
		a[1]: "..HH..HH..HH..HH",
		b[0]: "....HHHH....HHHH",
		b[1]: "........HHHHHHHH",
	}
	result := runScenario2(ci, scenario, "")

	if result["C"] != ".......H..HH.HHH" || result["S0"] != ".H.HH.H..H.HH.H." || result["S1"] != "..HH.HH.HH..H..H" {
		t.Error(result)
	}

}

func TestNbitConstant(t *testing.T) {
	ci := NewCircuit()

	c := ci.buildNbitConstant("test", 4, 5)
	ci.monitor(c[0], "C0")
	ci.monitor(c[1], "C1")
	ci.monitor(c[2], "C2")
	ci.monitor(c[3], "C3")

	dummy := ci.w()
	scenario := map[int]string{
		dummy: ".",
	}
	result := runScenario2(ci, scenario, "")

	if result["C0"]+result["C1"]+result["C2"]+result["C3"] != "H.H." {
		t.Error(result)
	}
}
