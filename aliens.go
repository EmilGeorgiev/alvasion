package alvasion

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
)

// AlienCommander is a commander of the aliens/soldiers. It is responsible for starting the invasion, distribute
// soldiers around the map, gives orders to which outgoing road the soldiers should continue.
// It has the map of the world and a list with all soldiers. At any moment it knows where are his soldiers and what is
// the current status of the invasion. It is main manager of the invasion.
type AlienCommander struct {
	WorldMap       WorldMap
	Soldiers       []*Alien
	TrappedAliens  int
	KilledSoldiers int
	Sitreps        chan Sitrep

	iterations    int
	IterationDone chan struct{}
	wg            sync.WaitGroup
}

// GenerateReportForInvasion generates a report about the state of the cities after the invasion. It only includes
// cities that haven't been destroyed. The method is used after the invasion is finished and commander can give his report.
func (ac *AlienCommander) GenerateReportForInvasion() string {
	buf := bytes.NewBufferString("")
	for _, city := range ac.WorldMap.Cities {
		if city.IsDestroyed {
			continue
		}
		row := city.Name
		for _, road := range city.OutgoingRoadsNames {
			if road == "" {
				continue
			}
			row += " " + road
		}
		buf.WriteString(row + "\n")
	}

	return buf.String()
}

// NewAlienCommander initialize and return a new AlienCommander.
func NewAlienCommander(wm WorldMap, aliens []*Alien, sitreps chan Sitrep) AlienCommander {
	return AlienCommander{
		WorldMap:       wm,
		Soldiers:       aliens,
		TrappedAliens:  0,
		KilledSoldiers: 0,
		Sitreps:        sitreps,
		IterationDone:  make(chan struct{}),
	}
}

// GiveOrdersToTheAlienIn gives orders to the alien in a specified city. If the city is destroyed or doesn't have an
// alien, it does nothing. If the alien is trapped or killed, it does nothing. If there are no available roads from the
// city, it increases the count of TrappedAliens. Otherwise, it randomly selects an available road and orders the alien
// to take it.
//
// This method is used in the begging of every invasion iteration. Every iteration starts sending orders to the soldiers.
func (ac *AlienCommander) GiveOrdersToTheAlienIn(c *City) {
	// if in the city there is an alien, then the commander will give orders to him
	if c.Alien == nil || c.IsDestroyed {
		c.Alien = nil
		return
	}

	if s := ac.Soldiers[c.Alien.ID]; s.Killed || s.Trapped {
		return
	}

	var availableRoads []chan Alien
	for _, r := range c.OutgoingRoads {
		if r == nil {
			continue
		}
		availableRoads = append(availableRoads, r)
	}

	if availableRoads == nil {
		ac.TrappedAliens++
		return
	}

	// The commander selects a random active outgoing road and orders the alien to take that road.
	i := rand.Intn(len(availableRoads))
	availableRoads[i] <- *c.Alien
	c.Alien = nil
}

// StartNextIteration starts the next iteration of the invasion. It orders all soldiers to invade, listens for situation
// reports, evaluates city destruction, waits for all evaluations to finish, and then evaluates road destruction.
//
// The method is used after the previous iteration is finished.
func (ac *AlienCommander) StartNextIteration() {
	wg := sync.WaitGroup{}

	// give orders to all soldiers
	for _, city := range ac.WorldMap.Cities {
		ac.GiveOrdersToTheAlienIn(city)
	}

	// prepare to listen for incoming situation reports about the evaluation of the invasion
	go ac.ListenForSitrep()

	// After issuing all orders, the commander can evaluate the consequences and assess the
	// damage inflicted upon the cities as a result of these commands.
	for _, city := range ac.WorldMap.Cities {
		EvaluateCityDestruction(city, &wg)
	}

	wg.Wait() // waiting the current iteration of the invasion to finish
	// send signal to notify that the iteration is finished. The commander should prepare the next iteration.
	ac.IterationDone <- struct{}{} // this will stop the sitrep listener, because no more reports will be sent
	// when a city is destroyed all roads leading out or in of the town also should be destroyed. It is important
	// to keep in mind that one road always connect two different cities
	wg2 := sync.WaitGroup{}
	for _, city := range ac.WorldMap.Cities {
		EvaluateRoadsDestruction(city, &wg2)
	}
	wg2.Wait()
}

// ListenForSitrep listens for situation reports from the soldiers. It updates the soldiers' statuses and the cities
// based on the reports. This listener is used during the current iteration. Soldiers send sitreps for the evaluation
// of the invasion. By this listener the commander know in every step where are his soldiers on the map and which cities
// are destroyed.
func (ac *AlienCommander) ListenForSitrep() {
	for {
		select {
		case report := <-ac.Sitreps:
			if !report.IsCityDestroyed {
				ac.WorldMap.Cities[report.CityName].Alien = ac.Soldiers[report.From]
				continue
			}

			ac.WorldMap.CleanCity(report.CityName)
			ac.Soldiers[report.From].Killed = true
			ac.KilledSoldiers++
			ac.WorldMap.Cities[report.CityName].Alien = nil
		default:
			// if there is no reports maybe the current iteration of
			// the invasion finished and all soldiers send their reports.
		}

		select {
		case <-ac.IterationDone:
			// current iteration is finished and all reports, from soldiers, are handled.
			// we can stop the listener until the next iteration begin.
			return
		default:
			// the current iteration of the invasion is in progress,
			// so we are continue to listen for reports from soldiers
		}
	}
}

func (ac *AlienCommander) StopInvasion() {

}

// DistributeAliens distributes the aliens randomly across the cities on the map. This method is used only once before
// first invasion of the iteration. In the beginning the commander can place one alien in a city. If the soldiers are more
// than cities, then part of the soldiers will not be distributed.
func (ac *AlienCommander) DistributeAliens() {
	var i int
	for _, city := range ac.WorldMap.Cities {
		if i >= len(ac.Soldiers) {
			break
		}
		city.Alien = ac.Soldiers[i]
		i++
	}
}

// StartInvasion  starts the invasion. It distributes the aliens, then starts a loop of iterations. After each iteration,
// it checks if the invasion should stop. If there have been 10,000 iterations or if there are fewer than two aliens left,
// it stops the invasion. This method is used only once when the commander start the invasion.
func (ac *AlienCommander) StartInvasion() {
	// as a fist step of invasion, the commander must distribute soldiers across cities on the map.
	ac.DistributeAliens()

	// the invasion is split in iterations. Every iteration has start, progress (number of steps), finish.
	// These iterations are repeated until the invasion finished. After every iteration finished, the commander
	// decide whether the invasion can continue to the next iteration or it should be interrupted.
	for {
		ac.StartNextIteration()
		ac.iterations++

		if ac.iterations >= 10000 {
			fmt.Println("Stop the invasion because the 10000 iterations ware made.")
			return
		}

		availableAliens := len(ac.Soldiers) - ac.KilledSoldiers - ac.TrappedAliens
		if availableAliens < 2 {
			fmt.Println("Stop the invasion because less then 2 aliens are available.")
			return
		}
	}
}

// Alien represent an alien soldier.
type Alien struct {
	// ID unique identification of the alien.
	ID int

	// Sitrep (Situation Report) is used by the alien to describe the current status of an ongoing mission to his commander.
	Sitreps chan Sitrep

	// Killed shows whether the alien is killed.
	Killed bool

	// Trapped shows whether the alien is trapped (When it is in a city without outgoing roads).
	Trapped bool
}

// Sitrep or situation report is used by aliens/soldiers to send reports about evaluations of the invasion.
type Sitrep struct {
	// From contains the id of the soldier/alien that sends the report
	From int

	// CityName is the name of the city from where this report is send
	CityName string

	// IsCityDestroyed give information to the commander whether the city is destroyed.
	IsCityDestroyed bool
}
