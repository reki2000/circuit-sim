package sim

import "strconv"

func setup() {
	gnd, vdd := w(), w()
	addModule(&Static{gnd, 0, "Gnd"})
	addModule(&Static{vdd, 1, "Vdd"})
}

func buildNandGate(name string) (in1, in2, out int) {
	in1, in2, out = w(), w(), w()
	out2 := w()
	addModule(&Mos{1, in1, out, true, name + ".p0"})
	addModule(&Mos{1, in2, out, true, name + ".p1"})
	addModule(&Mos{0, in2, out2, false, name + ".n1"})
	addModule(&Mos{out2, in1, out, false, name + ".n0"})
	return
}

func buildNotGate(name string) (in, out int) {
	in, out = w(), w()
	addModule(&Mos{1, in, out, true, name + ".p"})
	addModule(&Mos{0, in, out, false, name + ".n"})
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
