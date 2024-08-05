package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Line is a struct representing a line from the file
// with its corresponding number and text.
type Line struct {
	Text   string
	Number int64
}

// ReadLines opens a file and reads its lines one by one,
// sending them to a channel for processing.
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

// ValidateLines reads lines from a channel, validates the format,
// and splits them into parts for further processing.
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

var allPathsBetweenCities map[string]Path

// GenerateWorldMap creates a world map from the validated parts.
// It reads parts from a channel, and for each part creates a city
// with its roads, and adds it to the world map.
func GenerateWorldMap(parts <-chan []string) map[string]City {
	wolrdMap := map[string]City{}

	for p := range parts {
		var paths []Path
		for _, r := range p[1:] {
			rp := strings.Split(r, "=")
			outDirection := fmt.Sprintf("%s-%s", p[0], rp[1])
			inDirection := fmt.Sprintf("%s-%s", rp[1], p[0])

			path, ok := allPathsBetweenCities[outDirection]
			if ok {
				paths = append(paths, path)
				continue
			}

			ch1 := make(chan Alien)
			ch2 := make(chan Alien)

			p1 := Path{
				OutgoingDirection: ch1,
				IncomingDirection: ch2,
			}
			allPathsBetweenCities[outDirection] = p1
			p2 := Path{
				OutgoingDirection: ch2,
				IncomingDirection: ch1,
			}
			allPathsBetweenCities[inDirection] = p2
		}
		wolrdMap[p[0]] = NewCity(p[0], paths)
	}

	return wolrdMap
}
