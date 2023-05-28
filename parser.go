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
		Cities: map[string]*City{},
		Roads:  map[string]chan Alien{},
	}
	for p := range parts {
		city := &City{
			Name:          p[0],
			IncomingRoads: make([]chan Alien, 4),
			OutgoingRoads: make([]chan Alien, 4),
		}
		for _, r := range p[1:] {
			rp := strings.Split(r, "=")

			// 0 (north), 1 (south), 2 (east), 3 (west)
			if strings.ToLower(rp[0]) == "north" {
				outgoing := make(chan Alien, 1)
				city.OutgoingRoads[0] = outgoing
				//wm.Roads[road] = outgoing

				cityOnNorth, ok := wm.Cities[rp[1]]
				if ok {
					city.IncomingRoads[0] = cityOnNorth.OutgoingRoads[1]
					cityOnNorth.IncomingRoads[1] = city.OutgoingRoads[0]
				}
			} else if strings.ToLower(rp[0]) == "south" {
				outgoing := make(chan Alien, 1)
				city.OutgoingRoads[1] = outgoing
				//wm.Roads[road] = outgoing

				cityOnSouth, ok := wm.Cities[rp[1]]
				if ok {
					city.IncomingRoads[1] = cityOnSouth.OutgoingRoads[0]
					cityOnSouth.IncomingRoads[0] = city.OutgoingRoads[1]
				}
			} else if strings.ToLower(rp[0]) == "east" {
				outgoing := make(chan Alien, 1)
				city.OutgoingRoads[2] = outgoing
				//wm.Roads[road] = outgoing

				cityOnEast, ok := wm.Cities[rp[1]]
				if ok {
					city.IncomingRoads[2] = cityOnEast.OutgoingRoads[3]
					cityOnEast.IncomingRoads[3] = cityOnEast.OutgoingRoads[2]
				}
			} else if strings.ToLower(rp[0]) == "west" {
				outgoing := make(chan Alien, 1)
				city.OutgoingRoads[3] = outgoing
				//wm.Roads[road] = outgoing

				cityOnWest, ok := wm.Cities[rp[1]]
				if ok {
					city.IncomingRoads[3] = cityOnWest.OutgoingRoads[2]
					cityOnWest.IncomingRoads[2] = city.OutgoingRoads[3]
				}
			}
		}
		wm.Cities[city.Name] = city
	}
	return wm
}
