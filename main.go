package main

import "fmt"

var N, M int

type Grilla [][]string


func crearGrillaVacia(M int, N int) Grilla {
	var grilla Grilla
	grilla = make([][]string, M)
	for i := 0; i < M; i++ {
		grilla[i] = make([]string, N)
	}
	return grilla
}
func llenarGrilla(grilla Grilla, M int, N int) {
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			grilla[i][j] = "*"
		}
	}
}
func mostrarGrilla(grilla Grilla, M int, N int) {
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			fmt.Print(grilla[i][j])
		}
		fmt.Println("")
	}
}
func main() {
	M = 5
	N = 10
	var grilla = crearGrillaVacia(M, N)
	llenarGrilla(grilla, M, N)
	mostrarGrilla(grilla, M, N)
}
