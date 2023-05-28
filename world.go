package alvasion

import (
	"fmt"
	"sync"
)

type WorldMap struct {
	Cities map[string]*City
	Roads  map[string]chan Alien
}

func (wm WorldMap) CleanCity(name string) {
	c, ok := wm.Cities[name]
	if !ok || c.IsDestroyed {
		return
	}
	c.IsDestroyed = true

	// destroy and all roads leading out and in to the city
	for i, r := range c.OutgoingRoads {
		if r == nil {
			continue
		}
		// we can close the channel safely because
		close(r)
		c.OutgoingRoads[i] = nil
	}
}

type Road struct {
	Name string
	Ch   chan Alien
}

// City for simplicity we will add a convention that the in/out roads North, South, East,
// and West will be always in slace's indexes 0 (north), 1 (south), 2 (east), 3 (west).
type City struct {
	Name               string
	OutgoingRoads      []chan Alien
	IncomingRoads      []chan Alien
	IsDestroyed        bool
	Alien              Alien
	OutgoingRoadsNames []string
}

// EvaluateCityDestruction ..
func EvaluateCityDestruction(c *City, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		var aliens []Alien
		for _, road := range c.IncomingRoads {
			select {
			case alien := <-road:
				aliens = append(aliens, alien)
			default:
				// No alien is coming from this road in this invasion iteration
			}
		}

		if aliens == nil {
			return
		}

		if len(aliens) > 1 {
			msg := fmt.Sprintf("%s has been destroyed by alien ", c.Name)
			for _, a := range aliens {
				msg += fmt.Sprintf("%d%s", a.ID, " and alien ")
				a.Sitreps <- Sitrep{From: a.ID, CityName: c.Name, IsCityDestroyed: true}
			}
			fmt.Println(msg[:len(msg)-11] + "!")
			return
		}

		a := aliens[0]
		c.Alien = a
		a.Sitreps <- Sitrep{From: a.ID, CityName: c.Name, IsCityDestroyed: false}
	}()
}

// EvaluateRoadDestruction ...
func EvaluateRoadsDestruction(c *City, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i, road := range c.IncomingRoads {
			select {
			case _, ok := <-road:
				if !ok {
					// the road/channel is destroyed, and we can remove it from the
					// list with roads that leading in and out of the city
					c.IncomingRoads[i] = nil
					c.OutgoingRoads[i] = nil
					c.OutgoingRoadsNames[i] = ""
				}
			default:

			}
		}
	}()
}

var inOutMap = map[int]int{
	0: 1, // north-south
	1: 0, // south-north
	2: 3, // east-west
	3: 2, // west-east
}

//func generateWorldMap(fileName string) WorldMap {
//	file, err := os.Open(fileName)
//	if err != nil {
//		log.Fatalf("failed to open file: %s", err)
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	scanner.Split(bufio.ScanLines)
//
//	wm := WorldMap{
//		Cities: map[string]City{},
//		Roads:  map[string]chan Alien{},
//	}
//
//	for scanner.Scan() {
//		line := scanner.Text()
//		items := strings.Split(line, " ")
//		c := City{}
//		c.ID = items[0]
//
//		for _, road := range items[1:] {
//			r := strings.Split(road, "=")
//			if strings.ToLower(r[0]) == "west" {
//				// check whether the road already exist
//				//wm.Roads[]
//
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", "west", r[1])] = ch
//				c.West = ch
//			}
//
//			if strings.ToLower(r[0]) == "north" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.ID, "nort")] = ch
//				c.North = ch
//			}
//
//			if strings.ToLower(r[0]) == "east" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.ID, "east")] = ch
//				c.East = ch
//			}
//
//			if strings.ToLower(r[0]) == "south" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.ID, "south")] = ch
//				c.South = ch
//			}
//		}
//		wm.Cities[c.ID] = c
//	}
//	return wm
//}
//
//func createCity(line string) City {
//	return City{}
//}
