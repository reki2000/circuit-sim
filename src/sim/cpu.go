package sim

func (ci *Circuit) buildIP() (clk int, ip []int) {
	clk = ci.w()
	ip = []int{ci.w(), ci.w(), ci.w(), ci.w()}
	d, q, _clk := ci.buildNbitDFlipFlop("ip", 4)
	one := ci.buildNbitConstant("ip.1", 4, 1)
	a, b, _, s := ci.buildNbitAdder("ip.fa", 4)
	ci.bond(clk, _clk)
	ci.bondWires(one, b)
	ci.bondWires(q, a)
	ci.bondWires(s, d)
	ci.bondWires(q, ip)
	return
}
