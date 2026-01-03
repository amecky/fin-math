package math

import (
	"log"
	"math"
	"math/rand"
	"time"
)

// --------------------------------------
// K-Means Clustering
// --------------------------------------
func KMeans(m *Matrix, k int, iterations int) ([]float64, []int) {

	if k <= 0 {
		log.Fatal("k must be > 0")
	}

	rand.Seed(time.Now().UnixNano())

	// Randomly select initial centroids
	centroids := make([]float64, k)
	for i := 0; i < k; i++ {
		idx := rand.Intn(m.Rows)
		centroids[i] = m.DataRows[idx].Close()
	}

	assignments := make([]int, m.Rows)

	for iter := 0; iter < iterations; iter++ {

		// Assign points to closest centroid
		for i, v := range m.DataRows {
			minDist := math.Abs(v.Close() - centroids[0])
			minIdx := 0

			for c := 1; c < k; c++ {
				dist := math.Abs(v.Close() - centroids[c])
				if dist < minDist {
					minDist = dist
					minIdx = c
				}
			}

			assignments[i] = minIdx
		}

		// Recompute centroids
		counts := make([]int, k)
		sums := make([]float64, k)

		for i, c := range assignments {
			sums[c] += m.DataRows[i].Close()
			counts[c]++
		}

		for c := 0; c < k; c++ {
			if counts[c] > 0 {
				centroids[c] = sums[c] / float64(counts[c])
			}
		}
	}

	return centroids, assignments
}
