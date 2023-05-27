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

func ReadLines(fileName string, lines chan<- Line) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
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
	return nil
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

//func ParseLines(parts <-chan []string) {
//	wm := WorldMap{
//		Cities: map[string]City{},
//		Roads: map[string]chan Alien{},
//	}
//	for l := range lines {
//		parts := strings.Split(l.Text, " ")
//		if len(parts) < 2 {
//			fmt.Printf("Line number: %d has wrong format. Expect something like 'Foo west=Bar north=Baz' got: %s\n", l.Number, l.Text)
//			fmt.Println("Continue to the next line.")
//			continue
//		}
//		c := City{Name: parts[0]}
//
//		for i, road := range parts[1:] {
//			r := strings.Split(road, "=")
//			if len(r) != 2 {
//				fmt.Printf("On line %d the road number %d has wrong format. Expected something like 'west=Baz' got %s",l.Number, i, r)
//				fmt.Println("continue to the next road")
//				continue
//			}
//
//			if strings.ToLower(r[0]) == "west" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", "west", r[1])] = ch
//				c.West = ch
//			} else if strings.ToLower(r[0]) == "north" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "nort")] = ch
//				c.North = ch
//			} else if strings.ToLower(r[0]) == "east" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "east")] = ch
//				c.East = ch
//			} else if strings.ToLower(r[0]) == "south" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "south")] = ch
//				c.South = ch
//			} else {
//				fmt.Printf("On the line %d the road number %d has wrong direction. Expected 'west/north/east/south=Baz' got %s",l.Number, i, r[0])
//				fmt.Println("continue to the next road")
//				continue
//			}
//
//
//
//
//
//
//		}
//		wm.Cities[c.Name] = c
//	}
//}
