package main

import (
	"fmt"
	"math"
)

const COMPARE_THRESH = 3 //Exclusive

/* VISITED: 	Used in a substring
				 CANNOT be used in another substring
 * MARKED:		Perhaps under threshold, cannot initiate substring;
 				 CAN be used in another substring
 * UNVISITED:	Can initiate substring;
 				 CAN be used in another substring
*/
const UNVISITED = 0
const MARKED = 1
const VISITED = 2

func main() {

	s1 := []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")
	s2 := []byte("On a dark desert highway. I wanted to ride my bicycle.")

	fmt.Println(len(s1), len(s2))

	a, i := LCSubstr(s1, s2)
	print2d(a, i)
	fmt.Println(compareInter(a))

	// fmt.Println("Old", "Oldie")
	// print2d(LCSubstr("Old", "Oldie"))
}

func print2d(a [][]int, i int) {
	fmt.Println(i)
	for _, x := range a {
		for _, y := range x {

			if y == 0 {
				fmt.Print("- ")
			} else {
				fmt.Print(y, " ")
			}
		}
		fmt.Println()
	}
}

type Substring struct {
	value      []byte
	start      int //inclusive
	end        int //exclusive
	len        int
	fullString []byte
}

func compareInter(lcSubstrMatrix [][]int) []Substring {

	//Reduce matrix by turning the array of arrays into array of maxes
	var reduceSm []int
	for _, x := range lcSubstrMatrix {

		max, _ := max(x...)
		reduceSm = append(reduceSm, max)
	}

	//Preprocess so that vals under threshold are marked but not visited
	visited := make([]int, len(reduceSm))
	for i, _ := range visited {

		if reduceSm[i] < COMPARE_THRESH {
			visited[i] = MARKED
		}
	}

	//Until there are no more empty spaces, find comparisions (above certain threshold)
	var substrings []Substring
	for contains(UNVISITED, visited) != -1 {

		//Find curent max val
		max := 0
		maxI := -1
		for i, a := range reduceSm {
			if a > max && visited[i] == UNVISITED {
				max = a
				maxI = i
			}
		}

		//Create a substr using that max value
		substr := Substring{
			value:      nil, //reduceSm[maxI-max:maxI], //TODO figure out how to elegantly get this info
			start:      maxI - max,
			end:        maxI,
			len:        max,
			fullString: nil, //reduceSm,
		}

		fmt.Println(substr.end, visited)
		fmt.Println(reduceSm)

		//Confirm it is a valid substr that doesn't overlap
		checkVisit := visited[substr.start : substr.end+1]
		if contains(VISITED, checkVisit) != -1 {

			//If invalid, mark all substring as visited
			for i := substr.start; i < substr.end+1; i++ {
				visited[i] = int(math.Max(float64(visited[i]), float64(MARKED))) //Don't override if VISITED
			}
			continue
		}

		substrings = append(substrings, substr)

		//Visit indices used in this substr
		for i := substr.start; i < substr.end+1; i++ {
			visited[i] = VISITED
		}
	}

	return substrings
}

/* x is the one to check against, y is to check */
func LCSubstr(x []byte, y []byte) ([][]int, int) {
	m := len(x)
	n := len(y)

	lcSuff := make([][]int, 0)
	for i := 0; i < m+1; i++ {
		lcSuff = append(lcSuff, make([]int, n+1))
	}
	result := 0 //Store len of longest common substr

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {

			//1st row/col have no meaning
			if i == 0 || j == 0 {
				lcSuff[i][j] = 0

			} else if x[i-1] == y[j-1] {

				lcSuff[i][j] = lcSuff[i-1][j-1] + 1
				result = int(math.Max(float64(result), float64(lcSuff[i][j])))

			} else {

				lcSuff[i][j] = 0
			}
		}
	}

	return lcSuff, result
}
