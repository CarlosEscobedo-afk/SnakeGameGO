package main

var N, M int

type Grilla [][]string

/*
	for i := 0; i < 5; i++ {
		for j := 0; j < 7; j++ {
			fmt.Print(grilla[i][j])
		}
		fmt.Println("")
	}
*/

func crearGrillaVacia(M int, N int) Grilla {
	var grilla Grilla
	grilla = make([][]string, M)
	for i := 0; i < M; i++ {
		grilla[i] = make([]string, N)
	}
	return grilla
}

func main() {
	M = 10
	N = 9
	var grilla = crearGrillaVacia(M, N)
}
