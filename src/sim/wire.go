package sim

const numWires = 100

var wire [numWires]int
var maxWireID = -1

func w() int {
	if maxWireID >= len(wire)-1 {
		panic("too many wires")
	}
	maxWireID++
	return maxWireID
}

var net struct {
	from [numWires][]int
}

func bond(w1, w2 int) {
	//fmt.Printf(" bonding %d to %d\n", w1, w2)
	if net.from[w1] != nil {
		net.from[w1] = append(net.from[w1], w2)
	} else {
		net.from[w1] = []int{w2}
	}
}

func listBonded(w int) []int {
	return net.from[w]
}
