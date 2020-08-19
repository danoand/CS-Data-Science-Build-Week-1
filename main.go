package main

import (
	"fmt"
	"math"
)

func calcEuclidDist(row1, row2 []float64) float64 {
	var retDist float64

	// Calculate the sq rt of the sum of the squares
	//   of distances between two vectors
	for i := 0; i < len(row1); i++ {
		retDist = retDist + math.Pow((row1[i]-row2[i]), 2)
	}

	return math.Sqrt(retDist)
}

func main() {
	var dataset = [][]float64{
		{2.7810836, 2.550537003, 0.0},
		{1.465489372, 2.362125076, 0.0},
		{3.396561688, 4.400293529, 0.0},
		{1.38807019, 1.850220317, 0.0},
		{3.06407232, 3.005305973, 0.0},
		{7.627531214, 2.759262235, 1.0},
		{5.332441248, 2.088626775, 1.0},
		{6.922596716, 1.77106367, 1.0},
		{8.675418651, -0.242068655, 1.0},
		{7.673756466, 3.508563011, 1.0},
	}

	row0 := dataset[0]

	for j, rw := range dataset {
		distance := calcEuclidDist(row0, rw)
		fmt.Println(j+1, " | ", distance)
	}
}
