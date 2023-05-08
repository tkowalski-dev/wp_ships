package statystyki

import (
	"bufio"
	"fmt"
	"os"
)

var single *Statystyki

type Statystyki struct {
	Games []min_stat
}

type min_stat struct {
	trafienia int
	strzaly   int
}

func (m *min_stat) DodajTrafiony() {
	m.strzaly++
	m.trafienia++
}

func (m *min_stat) DodajPudlo() {
	m.strzaly++
}

func (m *min_stat) ObliczSkutecznosc() string {
	if m.trafienia == 0 {
		return "..."
	}
	return fmt.Sprintf("%.2f%", m.strzaly/m.trafienia*100.0)
}

func GetInstance() *Statystyki {
	if single == nil {
		single = new(Statystyki)
		single.Games = make([]min_stat, 0)
	}
	return single
}

func (s *Statystyki) GetNewGame() *min_stat {
	nowa := min_stat{}
	s.Games = append(s.Games, nowa)
	return &nowa
}

func (s *Statystyki) ObliczSkutecznosc() string {
	sumaTrafien := 0.0
	sumaStrzalow := 0.0
	for _, v := range s.Games {
		sumaTrafien += float64(v.trafienia)
		sumaStrzalow += float64(v.strzaly)
	}
	return fmt.Sprintf("%.2f%%", sumaTrafien/sumaStrzalow*100.0)
}

func (s *Statystyki) PokazStroneSkutecznosci() {
	fmt.Printf("### Statystyka skuteczności: ###\n")
	fmt.Printf("Ilosc gier: %v\n", len(s.Games))
	fmt.Printf("Ogólna skuteczność trafień: %v\n", s.ObliczSkutecznosc())
	fmt.Printf("\nSkuteczność w poszczególnych grach:\n")
	for i, v := range s.Games {
		fmt.Printf("Gra %v. Skutecznosc: %v\n", i, v.ObliczSkutecznosc())
	}

	fmt.Printf("\n-> Naciśnij enter, aby wrócić do menu...\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
