package main

func max(a ...int) (int, int) {

	max := 0
	idx := -1
	for i := 0; i < len(a); i++ {

		if a[i] > max {

			max = a[i]
			idx = i
		}
	}

	return max, idx
}

func contains(target int, a []int) int {

	for i, x := range a {

		if x == target {
			return i
		}
	}

	return -1
}
