package alvasion

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Line struct {
	Text   string
	Number int64
}

func ReadLines(fileName string, lines chan<- Line) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lineNumber int64
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		lines <- Line{Text: line, Number: lineNumber}
	}
	close(lines)
}

func ValidateLines(lines <-chan Line, p chan<- []string, errs chan<- error) {
LOOP:
	for l := range lines {
		parts := strings.Split(l.Text, " ")
		if len(parts) < 2 {
			errs <- fmt.Errorf("line number: %d has wrong format. A line should contains a city name and at least "+
				"one road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", l.Number, l.Text)
			continue
		}

		if len(parts) > 5 {
			errs <- fmt.Errorf("line number: %d has wrong format. A line should contains a city name and maximum "+
				"4 road that leading out of the city. Expect something like 'Foo west=Bar north=Baz' got: %s\n", l.Number, l.Text)
			continue
		}

		for i, road := range parts[1:] {
			r := strings.Split(road, "=")
			if len(r) != 2 {
				errs <- fmt.Errorf("on line %d the road number %d has wrong format. Expected something like 'west=Baz' got %s", l.Number, i+1, road)
				continue LOOP
			}

			rl := strings.ToLower(r[0])
			if rl != "west" && rl != "north" && rl != "east" && rl != "south" {
				errs <- fmt.Errorf("on the line %d the road number %d has wrong direction. Expected 'west/north/east/south' got %s", l.Number, i+1, r[0])
				continue LOOP
			}
		}

		p <- parts
	}
}

func GenerateWorldMap(parts <-chan []string) WorldMap {
	wm := WorldMap{
		Cities: map[string]City{},
		Roads:  map[string]chan Alien{},
	}
	for p := range parts {
		city := City{Name: p[0]}
		for _, r := range p[1:] {
			road := strings.ToLower(r)
			rp := strings.Split(road, "=")

			if strings.ToLower(rp[0]) == "west" {
				ch := make(chan Alien)
				a, ok := wm.Roads["east="+rp[1]]
				if ok {
					ch = a
				}
				city.Roads = append(city.Roads, ch)
			} else if strings.ToLower(rp[0]) == "north" {
				ch := make(chan Alien)
				a, ok := wm.Roads["south="+rp[1]]
				if ok {
					ch = a
				}
				wm.Roads[road] = ch
				city.Roads = append(city.Roads, ch)
			} else if strings.ToLower(rp[0]) == "east" {
				ch := make(chan Alien)
				a, ok := wm.Roads["west="+rp[1]]
				if ok {
					ch = a
				}
				wm.Roads[road] = ch
				city.Roads = append(city.Roads, ch)
			} else if strings.ToLower(rp[0]) == "south" {
				ch := make(chan Alien)
				a, ok := wm.Roads["north="+rp[1]]
				if ok {
					ch = a
				}
				wm.Roads[road] = ch
				city.Roads = append(city.Roads, ch)
			}
		}
		wm.Cities[city.Name] = city
	}
	return wm
}
