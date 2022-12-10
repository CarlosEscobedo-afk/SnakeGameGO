package main

/*
// REGLAS BÁSICAS:
Grilla de tamaño N×M.
Pueden haber 'p' serpientes.
Las serpientes no pueden transitar por una celda ocupada por otra serpiente.
Cada serpiente se desplaza de forma autónoma intentando llegar primero que las demás al alimento.
Cada vez que una serpiente alcance el alimento, ésta crecerá y, a la vez, deberá surgir una nueva porción de comida en otra celda (definida aleatoriamente).
Cada 'x' segundos cada serpiente se mueve una celda.
El juego debe continuar hasta que las serpientes no puedan seguir moviéndose.

// REGLAS CON CONCURRENCIA:
Las serpientes deben ser representadas por un Thread (liviano o pesado) independiente (con exclusion mutua).
Si una serpiente no puede moverse entra a deadlock permanente (o momentaneo).
El programa debe ser tan asíncrono y libre de barreras de sincronización como sea posible.

// GRILLA:
Defina un mecanismo que permita observar los diferentes estados de las serpientes en la grilla.
*/

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// MAIN------------------------------------------------------------------------------------------------------------------------------------------------
var ancho, largo int //Tamaño de la grilla.
var Dir byte         //Direccion de la serpiente.
var oldDir byte      //Direccion de la serpiente.
var posComida Coords

type MyGrilla [][]string //Grilla
var Score int

func (grilla MyGrilla) String() string {
	out := "█"
	for i := 0; i < len(grilla[0]); i++ {
		out += "▀"
	}
	out += "█\n"
	for i := 0; i < len(grilla); i++ {
		out += "█"
		for j := 0; j < len(grilla[i]); j++ {
			out += grilla[i][j]
		}
		out += "█\n"
	}
	out += "▀"
	for i := 0; i < len(grilla[0]); i++ {
		out += "▀"
	}
	out += "▀\n"
	return out
}

type Coords struct {
	X int
	Y int
}

type Snake struct {
	Cola []Coords
	lost bool
}

func (s *Snake) agregarCola(newLoc Coords) Coords {
	s.Cola = append(s.Cola, newLoc)
	return newLoc
}

func (s *Snake) quitarCola() Coords {
	temp := s.Cola[0]
	s.Cola = s.Cola[1:]
	return temp
}

//RETORNO DE CARRO GOLANG (hacer que la matriz se refresque bonito)

// FUNCIONES------------------------------------------------------------------------------------------------------------------------------------------------
func colocarComida(grilla MyGrilla) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	x, y := random.Intn(int(ancho-1)), random.Intn(int(largo-1))
	for grilla[x][y] != " " {
		x, y = random.Intn(int(ancho-1)), random.Intn(int(largo-1))
	}
	grilla[x][y] = "◈"
	posComida.X = x
	posComida.Y = y
}

func posicionInicialSerpientes(s *Snake, grilla MyGrilla) {
	s.lost = false
	time.Sleep(time.Millisecond * time.Duration(100))
	rand.Seed(time.Now().UnixNano())
	serpX, serpY := rand.Intn(int(ancho-1)), rand.Intn(int(largo-1))

	s.agregarCola(Coords{serpX, serpY + 1})
	grilla[serpX][serpY+1] = "□"
}

// 0 arriba;		1 abajo;	2 izquierda;	3 derecha
func verificar() {
	n, o := Dir, oldDir
	if (n == 0 && o == 1) || (n == 1 && o == 0) || (n == 3 && o == 2) || (n == 2 && o == 3) {
		Dir = oldDir
		fmt.Println("N:", n, "o:", o)
	}
}

func celdaSig(coord Coords, dir byte, grilla MyGrilla) (Coords, bool) {
	tempCoord := coord
	switch dir {
	case 0:
		tempCoord.X -= 1
		break
	case 2:
		tempCoord.Y -= 1
		break
	case 1:
		tempCoord.X += 1
		break
	case 3:
		tempCoord.Y += 1
		break
	}

	// Golpear la pared
	if tempCoord.X >= ancho || tempCoord.Y >= largo || tempCoord.X < 0 || tempCoord.Y < 0 {
		return tempCoord, false
	}
	// Comerse a si mismo
	if grilla[tempCoord.X][tempCoord.Y] != " " && grilla[tempCoord.X][tempCoord.Y] != "◈" {
		return tempCoord, false
	}
	return tempCoord, true
}

func actualizarGrilla(grilla MyGrilla, snake *Snake) {
	last := int(len(snake.Cola) - 1)
	newCell, cont := celdaSig(snake.Cola[last], Dir, grilla)

	if !cont {
		snake.lost = true
		return
	}

	snake.agregarCola(newCell)

	last = len(snake.Cola) - 1
	grilla[snake.Cola[last].X][snake.Cola[last].Y] = "□"

	if newCell == posComida {
		colocarComida(grilla)
	} else {
		grilla[snake.Cola[0].X][snake.Cola[0].Y] = " "
		snake.quitarCola()
	}
	oldDir = Dir
}

// MAIN------------------------------------------------------------------------------------------------------------------------------------------------
func main() {
	//Argumentos
	ancho1 := flag.Int("ancho", 0, "El ancho de la grilla.")
	largo1 := flag.Int("largo", 18, "El largo de la grilla.")
	vel1 := flag.Int("velocidad", 20, "La velocidad de la grilla.")
	flag.Parse()
	fmt.Println("Ancho = ", *ancho1)
	fmt.Println("Largo = ", *largo1)
	fmt.Println("Velocidad = ", *vel1)

	//Dimensiones de la grilla
	ancho, largo = 10, 25

	//Grilla Vacia
	var grilla MyGrilla
	grilla = make([][]string, ancho)
	for i := 0; i < int(ancho); i++ {
		grilla[i] = make([]string, largo)
	}

	//Inicializar Grilla
	for i := 0; i < int(ancho); i++ {
		for j := 0; j < int(largo); j++ {
			grilla[i][j] = string(" ")
		}
	}

	//Inicializar serpientes.
	//Serpiente 1
	var snake1 Snake
	ch1 := make(chan byte)

	go func(ch1 chan byte) {
		// LA POSICIÓN DE LA SERPIENTE NO DEBE CREAR UNA NUEVA GRILLA
		posicionInicialSerpientes(&snake1, grilla)
		//Movimiento serpiente
		for {
			last := int(len(snake1.Cola) - 1)

			time.Sleep(time.Millisecond * time.Duration(100))
			if posComida.X > snake1.Cola[last].X {
				time.Sleep(time.Millisecond * time.Duration(100))
				if Dir == 0 && snake1.Cola[last].Y-1 > 0 {
					oldDir = 0
					ch1 <- 2
				} else if Dir == 0 && snake1.Cola[last].Y-1 == 0 {
					oldDir = 0
					ch1 <- 3
				}
				ch1 <- 1
			} else if posComida.X < snake1.Cola[last].X {
				time.Sleep(time.Millisecond * time.Duration(100))
				ch1 <- 0
			} else if posComida.X == snake1.Cola[last].X {
				time.Sleep(time.Millisecond * time.Duration(100))
				if posComida.Y > snake1.Cola[last].Y {
					//Estos aun estan con pruebas
					//Condiciones para "evitar" Deadlock
					//if Dir == 2 && snake1.Cola[last].X-1 > 0 && grilla[snake1.Cola[last].X-1][snake1.Cola[last].Y] == "" {
					if Dir == 2 && snake1.Cola[last].X-1 > 0 {
						fmt.Println("<- to up")
						oldDir = 2
						ch1 <- 0
					} else if Dir == 2 && snake1.Cola[last].X-1 == 0 {
						fmt.Println("<- to down")
						oldDir = 2
						ch1 <- 1
					} else {
						ch1 <- 3
					}

				} else if posComida.Y < snake1.Cola[last].Y {
					if Dir == 3 && snake1.Cola[last].X-1 > 0 {
						fmt.Println("-> to up")
						oldDir = 3
						ch1 <- 0
					} else if Dir == 3 && snake1.Cola[last].X-1 == 0 {
						fmt.Println("-> to down")
						oldDir = 3
						ch1 <- 1
					} else {
						ch1 <- 2
					}
				}
			}
			//Deadlock
			if snake1.lost {
				for true {
				}
			}
		}
	}(ch1)

	//Colocar comida en la grilla
	colocarComida(grilla)
	fmt.Println(grilla)

	//Ejecución del juego
	for {
		//Direccion serpiente
		select {
		case stdin, _ := <-ch1:
			Dir = stdin
			verificar()
			actualizarGrilla(grilla, &snake1)
			time.Sleep(time.Millisecond * time.Duration(100))
		}
		/*
			case stdin, _ := <-ch2:
				Dir = stdin
				verificar()
				actualizarGrilla(grilla, &snake2)
				time.Sleep(time.Millisecond * time.Duration(80))
			}
		*/
		time.Sleep(time.Millisecond * time.Duration(200))
		fmt.Println(grilla) //Imprimir pasos de la grilla
		if snake1.lost == true {
			fmt.Println("Todas las serpientes han muerto.")
			break
		}

	}
}
