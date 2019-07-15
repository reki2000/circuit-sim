package sim

import "testing"

func TestRSLatch(t *testing.T) {
	setup()
	s, r, q, _ := buildRSLatch("test")

	scenario2 := map[int]string{
		r: "HHHHH.HHHHH",
		s: "H.HHHHHHH.H",
	}
	monitor(map[int]string{r: "r", s: "s", q: "q"})
	result := runScenario2(scenario2, "")

	if result["q"] != "HHHHH....HH" {
		t.Error()
	}
}

func TestGatedDLatch(t *testing.T) {
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
