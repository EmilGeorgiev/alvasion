package app

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type Randomizer interface {
	ChooseRoad([]chan Alien) chan Alien
}

// AlienCommander is a commander of the aliens/soldiers. It is responsible for starting the invasion, distribute
// soldiers around the map, gives orders to which outgoing road the soldiers should continue.
// It has the map of the world and a list with all soldiers. At any moment it knows where are his soldiers and what is
// the current status of the invasion. It is main manager of the invasion.
type AlienCommander struct {
	worldMap                []City
	soldiers                []Alien
	trappedAliens           atomic.Int64
	killedSoldiers          int
	sitreps                 chan Sitrep
	notifyDestroy           chan string
	randomizer              Randomizer
	writer                  io.Writer
	currentIterationReports []Sitrep
	maxNumberOfIteration    int

	iterations              int
	stopListeningForReports chan struct{}
	wait                    chan struct{}
	wg                      sync.WaitGroup
	mutex                   sync.Mutex
}

// GenerateReportForInvasion generates a report about the state of the cities after the invasion. It only includes
// cities that haven't been destroyed. The method is used after the invasion is finished and commander can give his report.
func (ac *AlienCommander) GenerateReportForInvasion() string {
	buf := bytes.NewBufferString("")
	for _, city := range ac.worldMap {
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
func NewAlienCommander(wm []City, aliens []Alien, r Randomizer, w io.Writer, n int) AlienCommander {
	return AlienCommander{
		worldMap:                wm,
		soldiers:                aliens,
		trappedAliens:           atomic.Int64{},
		killedSoldiers:          0,
		randomizer:              r,
		writer:                  w,
		maxNumberOfIteration:    n,
		stopListeningForReports: make(chan struct{}),
		wait:                    make(chan struct{}),
	}
}

// GiveOrdersToTheAlienIn gives orders to the alien in a specified city. If the city is destroyed or doesn't have an
// alien, it does nothing. If the alien is trapped or killed, it does nothing. If there are no available roads from the
// city, it increases the count of trappedAliens. Otherwise, it randomly selects an available road and orders the alien
// to take it.
//
// This method is used in the begging of every invasion iteration. Every iteration starts sending orders to the soldiers.
func (ac *AlienCommander) GiveOrdersToTheAlienIn(c City) {
	// if in the city there is an alien, then the commander will give orders to him
	alien := c.Alien
	if alien == nil || c.IsDestroyed {
		return
	}

	if s := ac.soldiers[alien.ID]; s.Killed || s.Trapped {
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
		ac.trappedAliens.Add(1)
		return
	}

	// The commander selects a random active outgoing road and orders the alien to take that road.
	road := ac.randomizer.ChooseRoad(availableRoads)
	road <- *alien
}

// StartIteration starts the next iteration of the invasion. It orders all soldiers to invade, listens for situation
// reports, evaluates city destruction, waits for all evaluations to finish, and then evaluates road destruction.
//
// The method is used after the previous iteration is finished.
func (ac *AlienCommander) StartIteration() {
	ac.sitreps = make(chan Sitrep, 1000)
	wg := sync.WaitGroup{}
	// give orders to all soldiers

	for i, city := range ac.worldMap {
		ac.GiveOrdersToTheAlienIn(city)
		city = city.SetAlien(nil).SetSitrep(ac.sitreps) // the city will become temporary without alien because the alien will continue to the next city.

		ac.worldMap[i] = city
	}
	// prepare to listen for incoming situation reports about the evaluation of the invasion

	go ac.ListenForSitrep()

	// After issuing all orders, the commander can evaluate the consequences and assess the
	// damage inflicted upon the cities as a result of these commands.
	for _, city := range ac.worldMap {
		if city.IsDestroyed {
			continue
		}
		wg.Add(1)
		go city.CheckForIncomingAliens(&wg)
	}
	wg.Wait() // waiting the current iteration of the invasion to finish
	// send signal to notify that the iteration is finished. The commander should prepare the next iteration.
	close(ac.sitreps) // this will stop the sitrep listener, because no more reports will be sent
	<-ac.wait         // wait all sitreps to be read from the commander
}

// ListenForSitrep listens for situation reports from the soldiers. It updates the soldiers' statuses and the cities
// based on the reports. This listeneEqualr is used during the current iteration. soldiers send sitreps for the evaluation
// of the invasion. By this listener the commander know in every step where are his soldiers on the map and which cities
// are destroyed.
func (ac *AlienCommander) ListenForSitrep() {
	for report := range ac.sitreps {
		if len(report.FromAliens) == 0 {
			continue
		}
		city := ac.worldMap[report.CityID]
		if city.IsDestroyed {
			continue
		}

		ac.currentIterationReports = append(ac.currentIterationReports, report)
	}
	ac.wait <- struct{}{} // indicate that all situation reports are read for the iteration.
}

func (ac *AlienCommander) validateSitrep(report Sitrep) bool {
	if len(report.FromAliens) == 0 {
		return false
	}

	city := ac.worldMap[report.CityID]
	return !city.IsDestroyed
}

// DistributeAliens distributes the aliens randomly across the cities on the map. This method is used only once before
// first invasion of the iteration. In the beginning the commander can place one alien in a city. If the soldiers are more
// than cities, then part of the soldiers will not be distributed.
func (ac *AlienCommander) DistributeAliens() {
	var i int
	for _, city := range ac.worldMap {
		if i >= len(ac.soldiers) {
			break
		}
		city.Alien = &ac.soldiers[i]
		ac.worldMap[city.ID] = city
		i++
	}
	ac.soldiers = ac.soldiers[:i]
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
		ac.StartIteration()
		ac.iterations++

		ac.applyChangesToWorldMapAfterIteration()
		if ac.iterations >= ac.maxNumberOfIteration {
			fmt.Println("Stop the invasion because the 10000 iterations ware made.")
			return
		}

		availableAliens := len(ac.soldiers) - ac.killedSoldiers - int(ac.trappedAliens.Load())
		if availableAliens < 2 {
			fmt.Printf("There is %d soldier left. Stop the invasion!\n", availableAliens)
			fmt.Println("Number of iterations: ", ac.iterations)
			return
		}
	}
}

func (ac *AlienCommander) applyChangesToWorldMapAfterIteration() {
	for _, sitrep := range ac.currentIterationReports {
		if len(sitrep.FromAliens) == 0 {
			city := ac.worldMap[sitrep.CityID]
			city = city.SetAlien(nil)
			ac.worldMap[city.ID] = city
		}

		if len(sitrep.FromAliens) == 1 {
			city := ac.worldMap[sitrep.CityID]
			city = city.SetAlien(&sitrep.FromAliens[0])
			ac.worldMap[city.ID] = city
			continue
		}

		if len(sitrep.FromAliens) > 1 {
			city := ac.worldMap[sitrep.CityID]
			city = city.Destroy()
			ac.worldMap[city.ID] = city
			msg := city.Name + " is destroyed from alien"
			for _, a := range sitrep.FromAliens {
				msg += fmt.Sprintf(" %d and alien", a.ID)
				ac.soldiers[a.ID].Killed = true
				ac.killedSoldiers++
			}
			fmt.Fprintln(ac.writer, msg[:len(msg)-10]+"!")
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(ac.worldMap))
	for _, city := range ac.worldMap {
		go func(c City) {
			c = c.CheckForDestroyedRoads(&wg)
			ac.worldMap[c.ID] = c
		}(city)
	}
	wg.Wait()
}

// Alien represent an alien soldier.
type Alien struct {
	// ID unique identification of the alien.
	ID int

	// Killed shows whether the alien is killed.
	Killed bool

	// Trapped shows whether the alien is trapped (When it is in a city without outgoing roads).
	Trapped bool
}

// Sitrep or situation report is used by aliens/soldiers to send reports about evaluations of the invasion.
type Sitrep struct {
	// From contains the id of the soldier/alien that sends the report
	FromAliens []Alien

	// CityName is the name of the city from where this report is send
	CityName string

	CityID int
}
