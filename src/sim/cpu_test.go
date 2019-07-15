package sim

import "testing"

func TestIP(t *testing.T) {
	ci := NewCircuit()

	_clk, ip := ci.buildIP()
	clk := ci.buildClock("CLK1", 2)
	ci.bond(clk, _clk)

	ci.monitor(clk, "CLK")
	ci.monitor(ip[0], "IP0")
	ci.monitor(ip[1], "IP1")
	ci.monitor(ip[2], "IP2")
	ci.monitor(ip[3], "IP3")

	dummy := ci.w()
	scenario := map[int]string{
		dummy: "............................................................................",
	}
	result := runScenario2(ci, scenario, "")

	if result["IP0"] != "HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH..HH.." {
		t.Error(result)
	}

}
