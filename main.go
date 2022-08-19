package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var matrixA, matrixB [][]int

// PopulateMatrix == populates matrix with random ints
func PopulateMatrix(size int) [][]int {
	rand.Seed(time.Now().UnixNano())
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrix[i] = append(matrix[i], rand.Intn(10))
		}
	}

	return matrix
}

// PrintMatrix == prints values stored iin matrix
func PrintMatrix(matrix [][]int, id string) {
	fmt.Println(id)

	//step == length of matrix multiplied by 100 to give goroutines time to execute
	step := time.Duration(len(matrix)*100) * time.Millisecond

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {

			fmt.Printf("%d ", matrix[i][j])
		}
		fmt.Println("\n")
		time.Sleep(step)
	}
}

// CreateEmpty == Fill Matrix with 0
func CreateEmpty(size int) [][]int {
	matrix := make([][]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrix[i] = append(matrix[i], 0)
		}
	}
	return matrix
}

// Track how long algorithms take
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start).Seconds()
	log.Printf("%s took %0.12f", name, elapsed)
}

// FoxAlgorithm == Fox's matrix multiplication algorithm
func FoxAlgorithm(wg *sync.WaitGroup, ch chan [][]int) {
	defer wg.Done()
	//defer timeTrack(time.Now(), "Fox")

	var q int //q == sq root of total matrix size/processes

	//bCastA == function which broadcasts aValues each stage to be multiplied with bValues
	bCastA := func(stage int, ch chan [][]int) {
		defer wg.Done()
		bCast := CreateEmpty(len(matrixA))

		for i := 0; i < len(matrixA); i++ {
			kBar := (i + stage) % q
			for j := 0; j < len(matrixA[i]); j++ {
				bCast[i][j] = matrixA[i][kBar]
			}
		}

		ch <- bCast
	}

	//bShift == function to shift up rows in matrix b after each stage
	bShift := func(ch chan [][]int) {
		defer wg.Done()

		ch <- matrixB
		shiftB := CreateEmpty(len(matrixB))
		keys := make([]int, len(matrixB))

		for i := 0; i < len(matrixB); i++ {
			for j := 0; j < len(matrixB[i]); j++ {
				shiftB[i][j] = matrixB[i][j]
			}
		}

		for i := len(matrixB) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(keys); i++ {
			matrixB[i] = shiftB[keys[i]]
		}
	}

	//c == Matrix Multiplication function which returns matrix c
	c := func() [][]int {
		matrixC := CreateEmpty(len(matrixA))

		//channels to store bValues & broadcasted aValues
		bs := make(chan [][]int)
		pr := make(chan [][]int)

		q = int(math.Sqrt(float64(len(matrixA) * len(matrixA[0]))))

		for stage := 0; stage < q; stage++ {
			wg.Add(1)
			go bCastA(stage, pr)

			wg.Add(1)
			go bShift(bs)

			wg.Add(1)
			go func() {
				defer wg.Done()

				a := <-pr
				b := <-bs
				for i := 0; i < len(matrixA); i++ {
					for j := 0; j < len(matrixA[i]); j++ {
						matrixC[i][j] = matrixC[i][j] + (a[i][j] * b[i][j])
					}
				}
			}()
		}

		return matrixC
	}

	ch <- c()
}

// CanonAlgorithm == Cannons Matrix Multiplication Algorithm
func CanonAlgorithm(wg *sync.WaitGroup, ch chan [][]int) {
	defer wg.Done()
	//defer timeTrack(time.Now(), "Cannon")

	//AlignMatrices == function which aligns Avalues left by i, and Bvalues up by j
	AlignMatrices := func() {

		//Align A values
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i < len(matrixA); i++ {
				row := make([]int, len(matrixA))
				for j := 0; j < len(matrixA[i]); j++ {
					row[j] = matrixA[i][j]
				}

				firstHalf := row[:i]
				secndHalf := row[i:]

				newRow := []int{}
				newRow = append(secndHalf, firstHalf...)

				for j := 0; j < len(matrixA); j++ {
					matrixA[i][j] = newRow[j]
				}
			}
		}()

		//Align B values
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i < len(matrixB); i++ {
				collumn := make([]int, len(matrixB))
				for j := 0; j < len(matrixB[i]); j++ {
					collumn[j] = matrixB[j][i]
				}

				firstHalf := collumn[:i]
				secndHalf := collumn[i:]

				newCollumn := []int{}
				newCollumn = append(secndHalf, firstHalf...)

				for j := 0; j < len(matrixB[i]); j++ {
					matrixB[j][i] = newCollumn[j]
				}
			}
		}()

		//For testing
		//PrintMatrix(matrixA, "new A")
		//PrintMatrix(matrixB, "new B")
	}

	//AShift == function to shift AValues left every phase
	AShift := func(ch chan [][]int) {
		defer wg.Done()

		shiftA := CreateEmpty(len(matrixA))
		keys := make([]int, len(matrixB))

		for i := len(matrixA) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(matrixA); i++ {
			for j := 0; j < len(matrixA[i]); j++ {
				shiftA[i][j] = matrixA[i][keys[j]]
			}
		}

		matrixA = shiftA
		ch <- shiftA
	}

	//BShift == function to shift BValues up every phase
	BShift := func(ch chan [][]int) {
		defer wg.Done()

		shiftB := CreateEmpty(len(matrixB))
		keys := make([]int, len(matrixB))

		for i := len(matrixB) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(keys); i++ {
			shiftB[i] = matrixB[keys[i]]
		}

		matrixB = shiftB
		ch <- shiftB
	}

	//c == Matrix Multiplication function which returns matrix c
	c := func() [][]int {
		matrixC := CreateEmpty(len(matrixA))

		phases := int(math.Sqrt(float64(len(matrixA) * len(matrixA[0]))))

		ach := make(chan [][]int)
		bch := make(chan [][]int)

		//Align
		AlignMatrices()
		for phase := 0; phase < phases; phase++ {

			wg.Add(1)
			go AShift(ach)

			wg.Add(1)
			go BShift(bch)

			wg.Add(1)
			go func() {
				defer wg.Done()

				a := <-ach
				b := <-bch

				for i := 0; i < len(matrixA); i++ {
					for j := 0; j < len(matrixA[i]); j++ {
						matrixC[i][j] = matrixC[i][j] + (a[i][j] * b[i][j])
					}
				}
			}()
		}
		return matrixC
	}

	ch <- c()
}

func main() {
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup

	//channels to store result matrix data
	fm := make(chan [][]int)
	cm := make(chan [][]int)

	matrixA = PopulateMatrix(5)
	PrintMatrix(matrixA, "A")

	matrixB = PopulateMatrix(5)
	PrintMatrix(matrixB, "B")

	wg.Add(1)
	go FoxAlgorithm(&wg, fm)
	foxMatrix := <-fm

	wg.Add(1)
	go CanonAlgorithm(&wg, cm)
	cannonMatrix := <-cm

	wg.Wait()

	PrintMatrix(foxMatrix, "C(Fox)")
	PrintMatrix(cannonMatrix, "C(Cannon)")

}
