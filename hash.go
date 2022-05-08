package main

import (
	"crypto/md5"
	"fmt"
	"math"
)

/*var uniset [][]uint8

func init() {

	var a, b, c uint8
	for a = 0; a < 128; a++ {
		for b = 0; b < 128; b++ {
			for c = 0; c < 128; c++ {
				uniset = append(uniset, []uint8{a, b, c})
			}
		}
	}
}/**/

func hash(data []uint8) []uint8 {

	hasher := md5.New()
	var output []uint8
	var currStr []uint8
	for i := 0; i < len(data); i++ {

		//fmt.Println("h", data[i])

		/* Break words on control characters, ie whitespace */
		if isControl(data[i]) {

			hasher.Write(currStr) //TODO not most concise way to do this. Should leverage pkg hash better
			output = append(output, hasher.Sum([]byte{})[0])

			currStr = []uint8{}
			hasher.Reset()
			continue
		}

		currStr = append(currStr, data[i])
	}
	output = append(output, hasher.Sum(currStr)[0]) //last word

	return output
}

func isControl(a uint8) bool {

	return a == 32 || a == 10 || a == 13
}

func quickCompare(a []uint8, b []uint8) {

	amap := make(map[uint8]int)
	for i := 0; i < len(a); i++ {
		amap[a[i]]++
	}

	bmap := make(map[uint8]int)
	for i := 0; i < len(b); i++ {
		bmap[a[i]]++
	}

	fmt.Println(amap)
	fmt.Println(bmap)

	var i uint8
	dataPoints := 0
	diffs := 0
	for i = 0; i < 128; i++ {

		d := int(math.Abs(float64(amap[i] - bmap[i])))
		p := int(math.Max(float64(amap[i]), float64(bmap[i])))
		dataPoints += p
		diffs += d
	}

	fmt.Printf("%d percent different. %d differences out of %d data points",
		int(float64(diffs)/float64(dataPoints)*100),
		diffs,
		dataPoints,
	)
}

func compare(a []uint8, b []uint8) {
	for i := 0; i < int(math.Max(float64(len(a)), float64(len(b)))); i++ {

		if a[i] == 97 {
			fmt.Print("-", "\t")
		} else {
			fmt.Print(a[i], "\t")
		}

		if b[i] == 97 {
			fmt.Println("-")
		} else {
			fmt.Println(b[i])
		}
	}
}
