package alvasion

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type WorldMap struct {
	Cities map[string]City
	Roads  map[string]chan Alien
}

func generateWorldMap(fileName string) WorldMap {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	wm := WorldMap{
		Cities: map[string]City{},
		Roads:  map[string]chan Alien{},
	}

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")
		c := City{}
		c.Name = items[0]

		for _, road := range items[1:] {
			r := strings.Split(road, "=")
			if strings.ToLower(r[0]) == "west" {
				// check whether the road already exist
				//wm.Roads[]

				ch := make(chan Alien)
				wm.Roads[fmt.Sprintf("%s-%s", "west", r[1])] = ch
				c.West = ch
			}

			if strings.ToLower(r[0]) == "north" {
				ch := make(chan Alien)
				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "nort")] = ch
				c.North = ch
			}

			if strings.ToLower(r[0]) == "east" {
				ch := make(chan Alien)
				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "east")] = ch
				c.East = ch
			}

			if strings.ToLower(r[0]) == "south" {
				ch := make(chan Alien)
				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "south")] = ch
				c.South = ch
			}
		}
		wm.Cities[c.Name] = c
	}
	return wm
}

func createCity(line string) City {
	return City{}
}
