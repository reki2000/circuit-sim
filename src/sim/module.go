package sim

import (
	"fmt"
	"strconv"
)

func (c *Circuit) buildClock(name string, duration int) (out int) {
	out = c.w()
	c.addDevice(&Clock{c.vdd, c.gnd, out, duration, name, false, duration})
	return
}

func (c *Circuit) buildNandGate(name string) (in1, in2, out int) {
	in1, in2, out = c.w(), c.w(), c.w()
	out2 := c.w()
	c.addDevice(&Mos{c.vdd, in1, out, true, name + ".p0"})
	c.addDevice(&Mos{c.vdd, in2, out, true, name + ".p1"})
	c.addDevice(&Mos{c.gnd, in2, out2, false, name + ".n1"})
	c.addDevice(&Mos{out2, in1, out, false, name + ".n0"})
	return
}

func (c *Circuit) buildNotGate(name string) (in, out int) {
	in, out = c.w(), c.w()
	c.addDevice(&Mos{c.vdd, in, out, true, name + ".p"})
	c.addDevice(&Mos{c.gnd, in, out, false, name + ".n"})
	return
}

func (c *Circuit) buildSomeNands(name string, countNand int) (in1, in2, out []int) {
	in1, in2, out = make([]int, countNand), make([]int, countNand), make([]int, countNand)
	for i := 0; i < countNand; i++ {
		myName := name + ".nand" + strconv.Itoa(i)
		in1[i], in2[i], out[i] = c.buildNandGate(myName)
	}
	return in1, in2, out
}

func (c *Circuit) bondNand(name string, in1, in2 int) (out int) {
	var _in1, _in2 int
	_in1, _in2, out = c.buildNandGate(name)
	c.bond(in1, _in1)
	c.bond(in2, _in2)
	return
}

func (c *Circuit) buildRSLatch(name string) (si, ri, q, qi int) {
	si, ri, q, qi = c.w(), c.w(), c.w(), c.w()
	in1, in2, out := c.buildSomeNands(name, 2)

	c.bond(si, in1[0])
	c.bond(ri, in1[1])
	c.bond(out[0], q)
	c.bond(out[0], in2[1])
	c.bond(out[1], qi)
	c.bond(out[1], in2[0])

	return
}

/*
 * Gated Delayed Latch is a tranparent latch. When CLK is On, Q reflects D. During CLK is Off, Q keeps previous Q
 * see https://ja.wikipedia.org/wiki/%E3%83%A9%E3%83%83%E3%83%81%E5%9B%9E%E8%B7%AF
 */
func (c *Circuit) buildGatedDLatch(name string) (clk, d, q int) {
	clk, d, q = c.w(), c.w(), c.w()

	in1, in2, out := c.buildSomeNands(name, 2)

	// d to nand0
	c.bond(d, in1[0])

	// clk to nand0/1
	c.bond(clk, in2[0])
	c.bond(clk, in2[1])

	// nand0 to nand1
	c.bond(out[0], in1[1])

	// first RS latch
	si, ri, _q, _ := c.buildRSLatch(name + ".rs")
	c.bond(out[0], si)
	c.bond(out[1], ri)

	c.bond(_q, q)
	return
}

func (c *Circuit) buildDFlipFlop(name string) (clk, d, q int) {
	clk, d, q = c.w(), c.w(), c.w()

	notin, notout := c.buildNotGate(name + ".not")
	clk1, d1, q1 := c.buildGatedDLatch(name + ".d1")
	clk2, d2, q2 := c.buildGatedDLatch(name + ".d2")

	// clk ot not
	c.bond(clk, notin)

	// not0 to nand0, nand1, not1
	c.bond(clk, clk2)
	c.bond(notout, clk1)
	c.bond(d, d1)
	c.bond(q1, d2)

	// d to nand0, not2
	c.bond(q2, q)

	return
}

func (c *Circuit) buildHalfAdder(name string) (a, b, ca, s int) {
	a, b, ca, s = c.w(), c.w(), c.w(), c.w()

	o0 := c.bondNand(name+".nand0", a, b)
	o1 := c.bondNand(name+".nand1", a, o0)
	o2 := c.bondNand(name+".nand2", b, o0)
	s = c.bondNand(name+".nand3", o1, o2)

	notin, notout := c.buildNotGate(name + ".not")
	c.bond(o0, notin)
	c.bond(notout, ca)

	return
}

func (c *Circuit) buildFullAdder(name string) (a, b, x, ca, s int) {
	a, b, x, ca, s = c.w(), c.w(), c.w(), c.w(), c.w()

	o0 := c.bondNand(name+".nand0", a, b)
	o1 := c.bondNand(name+".nand1", a, o0)
	o2 := c.bondNand(name+".nand2", b, o0)
	o3 := c.bondNand(name+".nand3", o1, o2)
	o4 := c.bondNand(name+".nand4", o3, x)
	o5 := c.bondNand(name+".nand5", o3, o4)
	o6 := c.bondNand(name+".nand6", o4, x)
	s = c.bondNand(name+".nand7", o5, o6)
	ca = c.bondNand(name+".nand8", o0, o4)

	return
}

func (c *Circuit) buildNbitAdder(name string, bits int) (a, b []int, ca int, s []int) {
	a, b, s = make([]int, bits), make([]int, bits), make([]int, bits)
	var lowCa int
	a[0], b[0], lowCa, s[0] = c.buildHalfAdder(name + ".ha")
	for i := 1; i < bits; i++ {
		var x int
		a[i], b[i], x, ca, s[i] = c.buildFullAdder(name + ".fa" + fmt.Sprintf("%d", i))
		c.bond(lowCa, x)
		lowCa = ca
	}
	return
}
