package main

func max(a ...int) (int, int) {

	max := 0
	idx := 0
	for i := 0; i < len(a); i++ {

		if a[i] > max {

			max = a[i]
			idx = i
		}
	}

	return max, idx
}
