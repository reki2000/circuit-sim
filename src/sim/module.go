package sim

import (
	"fmt"
	"strconv"
)

var gnd, vdd = w(), w()

func setup() {
	addDevice(&Static{gnd, 0, "Gnd"})
	addDevice(&Static{vdd, 1, "Vdd"})
}

func buildClock(name string, duration int) (out int) {
	out = w()
	addDevice(&Clock{vdd, gnd, out, duration, name, false, duration})
	return
}

func buildNandGate(name string) (in1, in2, out int) {
	in1, in2, out = w(), w(), w()
	out2 := w()
	addDevice(&Mos{vdd, in1, out, true, name + ".p0"})
	addDevice(&Mos{vdd, in2, out, true, name + ".p1"})
	addDevice(&Mos{gnd, in2, out2, false, name + ".n1"})
	addDevice(&Mos{out2, in1, out, false, name + ".n0"})
	return
}

func buildNotGate(name string) (in, out int) {
	in, out = w(), w()
	addDevice(&Mos{vdd, in, out, true, name + ".p"})
	addDevice(&Mos{gnd, in, out, false, name + ".n"})
	return
}

func buildSomeNands(name string, countNand int) (in1, in2, out []int) {
	in1, in2, out = make([]int, countNand), make([]int, countNand), make([]int, countNand)
	for i := 0; i < countNand; i++ {
		myName := name + ".nand" + strconv.Itoa(i)
		in1[i], in2[i], out[i] = buildNandGate(myName)
		monitor(map[int]string{
			in1[i]: myName + ".in1",
			in2[i]: myName + ".in2",
			out[i]: myName + ".out",
		})
	}
	return in1, in2, out
}

func bondNand(name string, in1, in2 int) (out int) {
	var _in1, _in2 int
	_in1, _in2, out = buildNandGate(name)
	bond(in1, _in1)
	bond(in2, _in2)
	monitor(map[int]string{
		in1:  name + ".in1",
		in2:  name + ".in2",
		_in1: name + ".in1",
		_in2: name + ".in2",
		out:  name + ".out",
	})
	return
}

func buildRSLatch(name string) (si, ri, q, qi int) {
	si, ri, q, qi = w(), w(), w(), w()
	in1, in2, out := buildSomeNands(name, 2)

	bond(si, in1[0])
	bond(ri, in1[1])
	bond(out[0], q)
	bond(out[0], in2[1])
	bond(out[1], qi)
	bond(out[1], in2[0])

	return
}

/*
 * Gated Delayed Latch is a tranparent latch. When CLK is On, Q reflects D. During CLK is Off, Q keeps previous Q
 * see https://ja.wikipedia.org/wiki/%E3%83%A9%E3%83%83%E3%83%81%E5%9B%9E%E8%B7%AF
 */
func buildGatedDLatch(name string) (clk, d, q int) {
	clk, d, q = w(), w(), w()

	in1, in2, out := buildSomeNands(name, 2)

	// d to nand0
	bond(d, in1[0])

	// clk to nand0/1
	bond(clk, in2[0])
	bond(clk, in2[1])

	// nand0 to nand1
	bond(out[0], in1[1])

	// first RS latch
	si, ri, _q, _ := buildRSLatch(name + ".rs")
	bond(out[0], si)
	bond(out[1], ri)

	bond(_q, q)
	return
}

func buildDFlipFlop(name string) (clk, d, q int) {
	clk, d, q = w(), w(), w()

	notin, notout := buildNotGate(name + ".not")
	clk1, d1, q1 := buildGatedDLatch(name + ".d1")
	clk2, d2, q2 := buildGatedDLatch(name + ".d2")

	// clk ot not
	bond(clk, notin)

	// not0 to nand0, nand1, not1
	bond(clk, clk2)
	bond(notout, clk1)
	bond(d, d1)
	bond(q1, d2)

	// d to nand0, not2
	bond(q2, q)

	monitor(map[int]string{clk1: "~CLK", q1: "Q1"})
	return
}

func buildHalfAdder(name string) (a, b, c, s int) {
	a, b, c, s = w(), w(), w(), w()

	o0 := bondNand(name+".nand0", a, b)
	o1 := bondNand(name+".nand1", a, o0)
	o2 := bondNand(name+".nand2", b, o0)
	s = bondNand(name+".nand3", o1, o2)

	notin, notout := buildNotGate(name + ".not")
	bond(o0, notin)
	bond(notout, c)

	monitor(map[int]string{
		a: name + ".A",
		b: name + ".B",
		c: name + ".C",
		s: name + ".S",
	})

	return
}

func buildFullAdder(name string) (a, b, x, c, s int) {
	a, b, x, c, s = w(), w(), w(), w(), w()

	o0 := bondNand(name+".nand0", a, b)
	o1 := bondNand(name+".nand1", a, o0)
	o2 := bondNand(name+".nand2", b, o0)
	o3 := bondNand(name+".nand3", o1, o2)
	o4 := bondNand(name+".nand4", o3, x)
	o5 := bondNand(name+".nand5", o3, o4)
	o6 := bondNand(name+".nand6", o4, x)
	s = bondNand(name+".nand7", o5, o6)
	c = bondNand(name+".nand8", o0, o4)

	monitor(map[int]string{
		a: name + ".A",
		b: name + ".B",
		x: name + ".X",
		c: name + ".c",
		s: name + ".S",
	})

	return
}

func buildNbitAdder(name string, bits int) (a, b []int, c int, s []int) {
	a, b, s = make([]int, bits), make([]int, bits), make([]int, bits)
	var lowC int
	a[0], b[0], lowC, s[0] = buildHalfAdder(name + ".ha")
	for i := 1; i < bits; i++ {
		var x int
		a[i], b[i], x, c, s[i] = buildFullAdder(name + ".fa" + fmt.Sprintf("%d", i))
		bond(lowC, x)
		lowC = c
	}
	return
}
