package alvasion

import (
	"sync"
)

type WorldMap struct {
	Cities map[string]*City
	//Roads  map[string]chan Alien
}

// Destroy marks the current city as destroyed and destroys all roads (channels) leading in or out of the city.
// If the city is already destroyed, the function returns the city unchanged. Otherwise, it marks the city as destroyed,
// makes all the incoming roads to be empty, and closes all outgoing roads.
//
// This function emulates the effect of an invasion on the city and its connected roads in the simulated world.
// When a city is destroyed, it's no longer accessible via any of the roads, so all the roads (represented by channels)
// are closed. This makes the city isolated from the rest of the cities in the world.
//
// Destroy operates on a copy of the City (as it does not have a pointer receiver),
// so remember to assign the result of the Destroy call to the original city if the changes need to be persisted.
//
// This method is used from the commander when 2 or more aliens visit the city, and they destroy it.
// AlienCommander use this method updated the world map that he has.
//
// Returns: A copy of the city with the applied changes.
func (c City) Destroy() City {
	if c.IsDestroyed {
		return c
	}
	c.IsDestroyed = true
	c.IncomingRoads = make([]chan Alien, 4) // destroy all incoming roads

	// destroy all outgoing roads.
	for i, r := range c.OutgoingRoads {
		if r == nil {
			continue
		}
		// by closing the channel we send event to the city that is on the
		// other side that this city is destroyed and this road can be used anymore.
		close(r)
		c.OutgoingRoads[i] = nil
		c.OutgoingRoadsNames[i] = ""
	}
	return c
}

// Road represents a path connecting two cities in the world.
// Each Road has a name, and a channel (Ch) that facilitates the movement of Aliens.
// Aliens use this channel to travel from one city to another during the invasion.
// In our model, every channel (road) has only one direction. This represents the fact
// that aliens can only travel in one direction along a road. For example, if there's a
// channel leading out from city A to city B, aliens can only use this channel to travel
// from A to B, not the other way around. To represent bidirectional travel between two cities,
// we use two separate channels. One for A to B, and another for B to A. This means that each pair
// of connected cities is linked by two channels: one for each direction of travel. This model
// ensures clarity in movement and road utilization, as each channel clearly defines the direction of
// travel between cities.
//
// Fields:
//   - Name: the unique identifier for the road.
//   - Ch: the channel used by Aliens to traverse this road. This channel essentially simulates the road,
//     allowing aliens to be "sent" from one city to another. If the channel is closed, this
//     implies that the road is destroyed and can no longer be used by aliens.
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
	Alien              *Alien
	OutgoingRoadsNames []string
}

// CheckForIncomingAliens checks all incoming roads for incoming alien soldiers.
// It runs concurrently, and when an alien is detected, the alien is added to a list of aliens.
//
// If the function detects one or more aliens arriving via the same road, it sends a situation report (Sitrep)
// containing the list of aliens and the name of the city. The Sitrep is sent through a channel of the first
// alien in the list. On the other side of the channel is the AlienCommander which determine what to do base on
// the information in the Sitrep.
//
// The function uses a WaitGroup to coordinate with other concurrently running tasks. The method should be use
// during the current iteration of the invasion after the commander gives orders to his soldiers.
//
// Parameters:
//   - wg: a pointer to a sync.WaitGroup used for managing concurrent tasks
//
// This function does not return any value.
func (c City) CheckForIncomingAliens(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		var aliens []Alien
		for _, road := range c.IncomingRoads {
			// If the road does not exist, skip this iteration
			if road == nil {
				continue
			}

			// Check for incoming alien
			select {
			case alien, ok := <-road:
				// If the road is destroyed or the alien is already killed (probably this should not happen), skip this iteration
				if !ok || alien.Killed {
					continue
				}
				aliens = append(aliens, alien)
			default:
			}
		}

		if len(aliens) == 0 {
			return
		}

		// If there are one or more aliens, an alien send a situation report to his commander.
		aliens[0].Sitreps <- Sitrep{FromAliens: aliens, CityName: c.Name}
	}()
}

// CheckForDestroyedRoads checks the status of each incoming road to the city.
// If a road (or channel) has been destroyed, it is removed from the city's incoming
// and outgoing road lists, and the corresponding entry in the city's OutgoingRoadsNames
// is cleared.
//
// This function is designed to run asynchronously in a goroutine, so it receives a WaitGroup
// as an argument to coordinate with other concurrently running tasks. Once it finishes
// its job, it calls Done on the WaitGroup to indicate completion.
//
// Use this method after all aliens finished with their job/movements during the iteration.
//
// The function returns the updated city.
//
// Parameters:
//   - wg: a pointer to a sync.WaitGroup used to coordinate with other concurrent tasks
//
// Returns:
//   - City: the updated city with the status of its roads checked and adjusted
func (c City) CheckForDestroyedRoads(wg *sync.WaitGroup) City {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i, road := range c.IncomingRoads {
			// Check each road to see if it has been destroyed
			select {
			case _, ok := <-road:
				if !ok {
					// If the road is destroyed, remove it from the incoming and outgoing roads
					// and clear the corresponding name from OutgoingRoadsNames
					c.IncomingRoads[i] = nil
					c.OutgoingRoads[i] = nil
					c.OutgoingRoadsNames[i] = ""
				}
				// If there is no destruction status available for this road, continue the loop
			default:

			}
		}
	}()
	// Return the updated city
	return c
}

//var inOutMap = map[int]int{
//	0: 1, // north-south
//	1: 0, // south-north
//	2: 3, // east-west
//	3: 2, // west-east
//}
//
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
