package sim

import "strconv"

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

	noti, noto := [3]int{}, [3]int{}
	for i := 0; i < 3; i++ {
		noti[i], noto[i] = buildNotGate(name + ".not" + string(i))
	}
	in1, in2, out := buildSomeNands(name, 8)

	// clk ot not0
	bond(clk, noti[0])

	// not0 to nand0, nand1, not1
	bond(noto[0], in2[0])
	bond(noto[0], in2[1])
	bond(noto[0], noti[1])

	// d to nand0, not2
	bond(d, in1[0])
	bond(d, noti[2])

	// not2 to nand1
	bond(noto[2], in1[1])

	// first RS latch
	bond(out[0], in1[2])
	bond(out[1], in2[3])
	bond(out[2], in1[3])
	bond(out[3], in2[2])

	// first RS latch to second
	bond(out[2], in1[4])
	bond(out[3], in1[5])

	// not1 to nadn4, nand5
	bond(noto[1], in2[4])
	bond(noto[1], in2[5])

	// final RS latch
	bond(out[4], in1[6])
	bond(out[5], in2[7])
	bond(out[6], in1[7])
	bond(out[7], in2[6])

	// output
	bond(out[6], q)
	return
}
