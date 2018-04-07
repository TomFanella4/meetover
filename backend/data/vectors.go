package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}
func foo(a *mat.VecDense) {
	fmt.Println(a.At(1, 0))
}
func flattenVector(rows int, vec mat.Matrix) float64 {
	res := 0.0
	for i := 0; i < rows; i++ {
		res += math.Abs(vec.At(i, 0))
	}
	return res
}
func vectorTest() {

	temp := mat.NewVecDense(3, nil)
	u := mat.NewVecDense(3, []float64{5, 2, 3})
	v := mat.NewVecDense(3, []float64{1, 9, 3})
	// s := []mat.Matrix{u, v}
	d := mat.Dot(u, v)
	temp.SubVec(u, v)
	fmt.Println("u dot v: ", d)
	matPrint(u)
	a := u.At(2, 0)
	fmt.Println(a)
	foo(temp)
	fmt.Println(flattenVector(3, temp))
}
