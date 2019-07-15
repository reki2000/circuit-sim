package sim

import "testing"

func TestRSLatch(t *testing.T) {
	resetRecord()
	setup()
	s, r, q, _ := buildRSLatch("test")

	scenario2 := map[int]string{
		r: "HHHHH.HHHHH",
		s: "H.HHHHHHH.H",
	}
	monitor(map[int]string{r: "r", s: "s", q: "q"})
	result := runScenario2(scenario2, "")

	if result["q"] != "HHHHH....HH" {
		t.Error(result)
	}
}

func TestGatedDLatch(t *testing.T) {
	resetRecord()
	setup()

	clk, d, q := buildGatedDLatch("d0")

	scenario2 := map[int]string{
		clk: "H...H.H....HHH.H.HHH.HH...",
		d:   "H.......HH.H......H..H....",
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	result := runScenario2(scenario2, "")

	if result["Q"] != "HHHH.......H......H..H...." {
		t.Error()
	}
}

func TestDFlipFlop(t *testing.T) {
	resetRecord()
	setup()

	clk, d, q := buildDFlipFlop("d")

	scenario2 := map[int]string{
		clk: ".H..H.H....HHH.H.HHH.HH...",
		d:   "HH......HHHH......H..H....",
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	result := runScenario2(scenario2, "")

	if result["Q"] != "HHHH.......HHHH..........." {
		t.Error()
	}
}

func TestClock(t *testing.T) {
	resetRecord()
	setup()

	clk := buildClock("CLK", 10)
	dummy := w()
	scenario := map[int]string{
		dummy: "............................................................",
	}
	monitor(map[int]string{clk: "CLK"})
	result := runScenario2(scenario, "")

	if result["CLK"] != "..........HHHHHHHHHH..........HHHHHHHHHH..........HHHHHHHHHH" {
		t.Error(result)
	}
}

func TestHalfAdder(t *testing.T) {
	resetRecord()
	setup()

	a, b, c, s := buildHalfAdder("test")

	scenario := map[int]string{
		a: "..HH",
		b: ".H.H",
	}
	monitor(map[int]string{a: "A", b: "B", c: "C", s: "S"})
	result := runScenario2(scenario, "")

	if result["C"] != "...H" || result["S"] != ".HH." {
		t.Error(result)
	}
}

func TestFullAdder(t *testing.T) {
	resetRecord()
	setup()

	a, b, x, c, s := buildFullAdder("test")

	scenario := map[int]string{
		a: "..HH..HH",
		b: ".H.H.H.H",
		x: "....HHHH",
	}
	monitor(map[int]string{a: "A", b: "B", x: "X", c: "C", s: "S"})
	result := runScenario2(scenario, "")

	if result["C"] != "...H.HHH" || result["S"] != ".HH.H..H" {
		t.Error(result)
	}
}

func TestNbitAdder(t *testing.T) {
	resetRecord()
	setup()

	a, b, c, s := buildNbitAdder("test", 2)

	scenario := map[int]string{
		a[0]: ".H.H.H.H.H.H.H.H",
		a[1]: "..HH..HH..HH..HH",
		b[0]: "....HHHH....HHHH",
		b[1]: "........HHHHHHHH",
	}
	monitor(map[int]string{a[0]: "A0", a[1]: "A1", b[0]: "B0", b[1]: "B1", c: "C", s[0]: "S0", s[1]: "S1"})
	result := runScenario2(scenario, "")

	if result["C"] != ".......H..HH.HHH" || result["S0"] != ".H.HH.H..H.HH.H." || result["S1"] != "..HH.HH.HH..H..H" {
		t.Error(result)
	}

}

func TestNbitAdder2(t *testing.T) {
	resetRecord()
	setup()

	a, b, _, _ := buildNbitAdder("test", 2)

	scenario := map[int]string{
		a[0]: ".",
		a[1]: ".",
		b[0]: ".",
		b[1]: ".",
	}
	result := runScenario2(scenario, "-")

	if result["C"] != "...H.HHH" || result["S0"] != ".HH.H..H" {
		t.Error(result)
	}

}
