package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NFICHAS       = 4
	MAXOBSTACULOS = 10
)

type Ficha struct {
	id       int
	color    string
	posicion int
	estado   int //0 si zona blanca //1 si toco zono obstaculo y recien va perder un turno //si esta en zona obstaculo pero ya perdio un turno
	meta     bool
}
type Lanzamiento struct {
	dadoA   int
	dadoB   int
	avanzar bool
}

func getNumberPlayers(n *int) {
	for *n > 4 || *n < 1 {
		fmt.Print("Número de jugadores: ")
		fmt.Scanf("%d\n", &*n)
	}
}

func loadGame(fichas *[]Ficha, tabla *[40]int, nPlayers int, colors []string, positions []int) {
	var contador int

	for contador < MAXOBSTACULOS {
		number := rand.Intn(40)
		found := false
		for _, v := range positions {
			if number == v {
				found = true
				break
			}
		}
		if !found {
			contador++
			(*tabla)[number] = -1
		}

	}

	for i := 0; i < nPlayers; i++ {
		for j := 0; j < NFICHAS; j++ {
			ficha := Ficha{
				id:       j + 1,
				color:    colors[i],
				posicion: positions[i],
				meta:     false,
			}
			*fichas = append(*fichas, ficha)
		}
	}
}

func lanzarDados() Lanzamiento {
	valor := rand.Intn(2)
	tiro := Lanzamiento{
		dadoA:   rand.Intn(6) + 1,
		dadoB:   rand.Intn(6) + 1,
		avanzar: valor == 1,
	}
	return tiro
}

func pierdeTurno(tabla [40]int, fichas *[]Ficha, n int) bool {

	for i := 0; i < 4; i++ {
		for ind, valor := range tabla {
			if valor == -1 && ind == (*fichas)[i+n].posicion {
				// fmt.Println("COINCEDENCIA")
				(*fichas)[i+n].estado += 1
				if (*fichas)[i+n].estado > 2 {
					(*fichas)[i+n].estado = 2
				}
			}
			if (*fichas)[i+n].estado == 2 && valor == 0 && ind == (*fichas)[i+n].posicion {
				(*fichas)[i+n].estado = 0
			}
		}
	}
	// fmt.Println((*fichas)[n : n+4])

	for i := 0; i < 4; i++ {
		if (*fichas)[i+n].estado == 1 {
			return true
		}
	}
	return false
}

func turnoJugador(color string, wg *sync.WaitGroup, miTurno chan bool, ficha1 chan bool, ficha2 chan bool, ficha3 chan bool, ficha4 chan bool, fichas *[]Ficha, tabla [40]int) {
	defer wg.Done()
	var finJuego bool
	var ind int
	var n int

	if color == "red" {
		n = 0
	}
	if color == "green" {
		n = 4
	}
	if color == "blue" {
		n = 8
	}
	if color == "yellow" {
		n = 12
	}

	for !finJuego {
		time.Sleep(time.Millisecond * 100)
		miTurno <- true
		fmt.Printf("TURNO JUGADOR %s \n", color)
		if !pierdeTurno(tabla, *&fichas, n) {
			fmt.Println("JUGANDO...")
			var tiro Lanzamiento = lanzarDados()
			go func() {
				if (*fichas)[n].meta == false && (*fichas)[n].posicion != 39 {
					ficha1 <- true
				}
			}()
			go func() {
				if (*fichas)[n+1].meta == false && (*fichas)[n+1].posicion != 39 {
					ficha2 <- true
				}
			}()
			go func() {
				if (*fichas)[n+2].meta == false && (*fichas)[n+2].posicion != 39 {
					ficha3 <- true
				}
			}()
			go func() {
				if (*fichas)[n+3].meta == false && (*fichas)[n+2].posicion != 39 {
					ficha4 <- true
				}
			}()

			select {
			case <-ficha1:
				fmt.Printf("(JUEGA FICHA 1)\n")
				ind = n

			case <-ficha2:
				fmt.Printf("(JUEGA FICHA 2)\n")
				ind = n + 1

			case <-ficha3:
				fmt.Printf("(JUEGA FICHA 3)\n")
				ind = n + 2

			case <-ficha4:
				fmt.Printf("(JUEGA FICHA 4)\n")
				ind = n + 3
			}

			go func() {
				for {
					select {
					case <-ficha1:
					case <-ficha2:
					case <-ficha3:
					case <-ficha4:
						// Descartar elementos del canal
					default:
						// El canal está vacío
						return
					}
				}
			}()

			if tiro.avanzar {
				fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA+tiro.dadoB)
				(*fichas)[ind].posicion += tiro.dadoA + tiro.dadoB
				if (*fichas)[ind].posicion > 39 {
					(*fichas)[ind].posicion = 39 - ((*fichas)[ind].posicion - 39)
				}
			} else {
				fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA-tiro.dadoB)
				(*fichas)[ind].posicion += tiro.dadoA - tiro.dadoB
				if (*fichas)[ind].posicion < 0 {
					(*fichas)[ind].posicion = 0
				}
			}
			fmt.Println("POSCION ACTUAL DE LA FICHA: ", (*fichas)[ind].posicion)

			for i := n; i < n+4; i++ {
				if (*fichas)[ind].posicion == 39 {
					(*fichas)[ind].meta = true
				}
			}
			fichasCompletadas := 0
			for _, f := range (*fichas)[n : n+4] {
				if f.meta == true {
					fichasCompletadas++
				}
			}

			if fichasCompletadas == 4 {
				fmt.Println("***ACABA DE GANAR EL COLOR: ", color)
				finJuego = true
			}

		} else {
			fmt.Println("ESTE JUGADOR PERDIO SU TURNO")
		}
		fmt.Println("------------------------")

		select {
		case <-miTurno:
			// El siguiente jugador está listo
		default:
			// No hay jugadores listos, así que espero
		}
	}
}

func main() {
	var wg sync.WaitGroup
	var tabla [40]int

	fichas := []Ficha{}
	colors := []string{"red", "green", "blue", "yellow"}
	positions := []int{0, 0, 0, 0, 39}

	var nPlayers int

	getNumberPlayers(&nPlayers)
	loadGame(&fichas, &tabla, nPlayers, colors, positions)
	fmt.Println(tabla)

	miTurno := make(chan bool, 1)
	// suTurno := make(chan bool, 1)
	// suTurno <- true // Inicialmente, el primer jugador puede avanzar
	chFichas := make([]chan bool, nPlayers*NFICHAS)

	for i := range chFichas {
		chFichas[i] = make(chan bool)
	}

	for ind, c := range colors[:nPlayers] {
		wg.Add(1)
		go turnoJugador(c, &wg, miTurno, chFichas[ind*NFICHAS], chFichas[ind*NFICHAS+1], chFichas[ind*NFICHAS+2], chFichas[ind*NFICHAS+3], &fichas, tabla)
	}

	wg.Wait()
	close(miTurno)
	// close(suTurno)

	fmt.Println("!!!JUEGO FINALIZADO TODOS LOS JUGADORES LLEGARON A LA META !!!")
	for _, v := range fichas {
		fmt.Println(v)
	}

}
