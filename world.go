package alvasion

import (
	"fmt"
	"sync"
)

type WorldMap struct {
	Cities map[string]City
	Roads  map[string]chan Alien
}

func (wm WorldMap) CleanCity(name string) {
	c, ok := wm.Cities[name]
	if !ok {
		return
	}
	c.IsDestroyed = true

	// destroy and all roads leading out and in to the city
	for _, r := range c.Roads {
		// we can close the channel safely because
		close(r)
	}
}

type City struct {
	Name        string
	Roads       []chan Alien
	IsDestroyed bool
	Alien       Alien
}

// EvaluateCityDestruction ..
func EvaluateCityDestruction(c *City, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		var aliens []Alien
		var destroyedRoads []int
		for i, road := range c.Roads {
			select {
			case alien, ok := <-road:
				if !ok {
					// the road/channel is destroyed, and we can remove it from the
					// list with roads that leading in and out of the city
					destroyedRoads = append(destroyedRoads, i)
					continue
				}
				aliens = append(aliens, alien)
			default:
				// No alien is coming from this road in this invasion iteration
			}
		}

		// remove destroyed roads. They will be no longer used by anyone.
		for j := len(destroyedRoads) - 1; j >= 0; j-- {
			i := destroyedRoads[j]
			c.Roads = append(c.Roads[:i], c.Roads[i+1:]...)
		}

		//select {
		//case alien, ok := <-c.West:
		//	if !ok {
		//		// channel is closed and is not possible an alien to come from it
		//		break
		//	}
		//	aliens = append(aliens, alien)
		//default:
		//	// No alien is coming from the west road in this invasion iteration
		//}
		//
		//select {
		//case alien, ok := <-c.South:
		//	if !ok {
		//		// channel is closed and is not possible an alien to come from it
		//		break
		//	}
		//	aliens = append(aliens, alien)
		//default:
		//	// No alien is coming from the south road in this invasion iteration
		//}
		//
		//select {
		//case alien, ok := <-c.North:
		//	if !ok {
		//		// channel is closed and is not possible an alien to come from it
		//		break
		//	}
		//	aliens = append(aliens, alien)
		//default:
		//	// No alien is coming from the north road in this invasion iteration
		//}

		if aliens == nil {
			fmt.Println("No Aliens")
			return
		}

		if len(aliens) > 1 {
			msg := fmt.Sprintf("%s has been destroyed by alien ", c.Name)
			for _, a := range aliens {
				msg += a.Name + " and alien "
				a.Sitreps <- Sitrep{From: a.Name, CityName: c.Name, IsCityDestroyed: true}
			}
			fmt.Println(msg[:len(msg)-11] + "!")
		}

		a := aliens[0]
		a.Sitreps <- Sitrep{From: a.Name, CityName: c.Name, IsCityDestroyed: false}
	}()
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
//		c.Name = items[0]
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
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "nort")] = ch
//				c.North = ch
//			}
//
//			if strings.ToLower(r[0]) == "east" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "east")] = ch
//				c.East = ch
//			}
//
//			if strings.ToLower(r[0]) == "south" {
//				ch := make(chan Alien)
//				wm.Roads[fmt.Sprintf("%s-%s", c.Name, "south")] = ch
//				c.South = ch
//			}
//		}
//		wm.Cities[c.Name] = c
//	}
//	return wm
//}
//
//func createCity(line string) City {
//	return City{}
//}
