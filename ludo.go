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
