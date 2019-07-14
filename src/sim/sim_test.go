package sim

import "testing"

func TestRSLatch(t *testing.T) {
	setup()
	s, r, q, _ := buildRSLatch("test")
	scenario := map[int][]int{
		r: {1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1},
		s: {1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
	}
	monitor(map[int]string{r: "r", s: "s", q: "q"})
	result := runScenario(scenario, "")

	if result["q"] != "HHHHH____HH" {
		t.Error()
	}
}

func TestGatedDLatch(t *testing.T) {
	setup()

	clk, d, q := buildGatedDLatch("d0")

	scenario := map[int][]int{
		clk: {1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0},
		d:   {1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	result := runScenario(scenario, "")

	if result["Q"] != "HHHH_______H______H__H____" {
		t.Error()
	}
}

func TestGatedDLatch2(t *testing.T) {
	setup()

	clk, d, q := buildGatedDLatch("d0")

	// wireID: { value series for each t ...}
	scenario := map[int][]int{
		clk: {1, 1},
		d:   {1, 0},
	}

	monitor(map[int]string{clk: "CLK", d: "D", q: "Q"})
	result := runScenario(scenario, "d0.nand3.n")

	if result["Q"] != "H_" {
		t.Error()
	}
}
