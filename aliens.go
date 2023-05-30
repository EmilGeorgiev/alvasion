package alvasion

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
)

// AlienCommander is a commander of the aliens/soldiers. It is responsible for starting the invasion, distribute
// soldiers around the map, gives orders to which outgoing road the soldiers should continue.
// It has the map of the world and a list with all soldiers. At any moment it knows where are his soldiers and what is
// the current status of the invasion. It is main manager of the invasion.
type AlienCommander struct {
	WorldMap       map[string]City
	Soldiers       []*Alien
	TrappedAliens  atomic.Int64
	KilledSoldiers int
	Sitreps        chan Sitrep
	NotifyDestroy  chan string

	iterations    int
	IterationDone chan struct{}
	wait          chan struct{}
	wg            sync.WaitGroup
}

// GenerateReportForInvasion generates a report about the state of the cities after the invasion. It only includes
// cities that haven't been destroyed. The method is used after the invasion is finished and commander can give his report.
func (ac *AlienCommander) GenerateReportForInvasion() string {
	buf := bytes.NewBufferString("")
	for _, city := range ac.WorldMap {
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

func (ac *AlienCommander) GenerateReportForInvasion2() map[string][]string {
	var keys []string
	for k, city := range ac.WorldMap {
		if city.IsDestroyed {
			continue
		}
		keys = append(keys, k)

	}
	sort.Strings(keys)

	result := map[string][]string{}
	for _, name := range keys {
		city := ac.WorldMap[name]
		for _, road := range city.OutgoingRoadsNames {
			if road == "" {
				continue
			}
			result[name] = append(result[name], road)
		}
	}

	return result
}

func (ac *AlienCommander) SetNotifyDestroy(n chan string) {
	ac.NotifyDestroy = n
}

// NewAlienCommander initialize and return a new AlienCommander.
func NewAlienCommander(wm map[string]City, aliens []*Alien, sitreps chan Sitrep) AlienCommander {
	return AlienCommander{
		WorldMap:       wm,
		Soldiers:       aliens,
		TrappedAliens:  atomic.Int64{},
		KilledSoldiers: 0,
		Sitreps:        sitreps,
		IterationDone:  make(chan struct{}),
		wait:           make(chan struct{}),
	}
}

// GiveOrdersToTheAlienIn gives orders to the alien in a specified city. If the city is destroyed or doesn't have an
// alien, it does nothing. If the alien is trapped or killed, it does nothing. If there are no available roads from the
// city, it increases the count of TrappedAliens. Otherwise, it randomly selects an available road and orders the alien
// to take it.
//
// This method is used in the begging of every invasion iteration. Every iteration starts sending orders to the soldiers.
func (ac *AlienCommander) GiveOrdersToTheAlienIn(c City) {
	// if in the city there is an alien, then the commander will give orders to him
	if c.Alien == nil || c.IsDestroyed {
		//c.Alien = nil
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
		ac.TrappedAliens.Add(1)
		return
	}

	// The commander selects a random active outgoing road and orders the alien to take that road.
	i := rand.Intn(len(availableRoads))
	availableRoads[i] <- *c.Alien
	//c.Alien = nil
}

// StartIteration starts the next iteration of the invasion. It orders all soldiers to invade, listens for situation
// reports, evaluates city destruction, waits for all evaluations to finish, and then evaluates road destruction.
//
// The method is used after the previous iteration is finished.
func (ac *AlienCommander) StartIteration() {
	wg := sync.WaitGroup{}
	// give orders to all soldiers
	var cities []City
	for k, city := range ac.WorldMap {
		ac.GiveOrdersToTheAlienIn(city)
		city.Alien = nil
		ac.WorldMap[k] = city
		cities = append(cities, city)
	}

	// prepare to listen for incoming situation reports about the evaluation of the invasion
	go ac.ListenForSitrep()

	// After issuing all orders, the commander can evaluate the consequences and assess the
	// damage inflicted upon the cities as a result of these commands.
	for _, city := range cities {
		city.CheckForIncomingAliens(&wg)
	}
	wg.Wait() // waiting the current iteration of the invasion to finish

	// send signal to notify that the iteration is finished. The commander should prepare the next iteration.
	ac.IterationDone <- struct{}{} // this will stop the sitrep listener, because no more reports will be sent
	<-ac.wait

	// when a city is destroyed all roads leading out or in of the town also should be destroyed. It is important
	// to keep in mind that one road always connect two different cities
	//wg.Add(len(ac.WorldMap))
	for _, city := range ac.WorldMap {
		newC := city.CheckForDestroyedRoads()
		ac.WorldMap[newC.Name] = newC
	}
}

// ListenForSitrep listens for situation reports from the soldiers. It updates the soldiers' statuses and the cities
// based on the reports. This listener is used during the current iteration. Soldiers send sitreps for the evaluation
// of the invasion. By this listener the commander know in every step where are his soldiers on the map and which cities
// are destroyed.
func (ac *AlienCommander) ListenForSitrep() {
	for {
		select {
		case report := <-ac.Sitreps:
			if !ac.validateSitrep(report) {
				continue
			}

			city := ac.WorldMap[report.CityName]
			if len(report.FromAliens) == 1 {
				city.Alien = &report.FromAliens[0]
				ac.WorldMap[report.CityName] = city
				continue
			}

			destroyedCity := city.Destroy()
			ac.WorldMap[destroyedCity.Name] = destroyedCity
			msg := report.CityName + " is destroyed from alien"
			for _, a := range report.FromAliens {
				msg += fmt.Sprintf(" %d and alien", a.ID)
				ac.Soldiers[a.ID].Killed = true
				ac.KilledSoldiers++
			}
			fmt.Println(msg[:len(msg)-10] + "!")
			if ac.NotifyDestroy != nil {
				ac.NotifyDestroy <- report.CityName
			}
		default:
			// if there is no reports maybe the current iteration of
			// the invasion finished and all soldiers send their reports.
		}

		select {
		case <-ac.IterationDone:
			// current iteration is finished and all reports, from soldiers, are handled.
			// we can stop the listener until the next iteration begin.
			ac.wait <- struct{}{}
			return
		default:
			// the current iteration of the invasion is in progress,
			// so we are continue to listen for reports from soldiers
		}
	}
}

func (ac *AlienCommander) StopInvasion() {

}

func (ac *AlienCommander) validateSitrep(report Sitrep) bool {
	if len(report.FromAliens) == 0 {
		return false
	}

	city, ok := ac.WorldMap[report.CityName]
	if !ok {
		return false
	}

	return !city.IsDestroyed
}

// DistributeAliens distributes the aliens randomly across the cities on the map. This method is used only once before
// first invasion of the iteration. In the beginning the commander can place one alien in a city. If the soldiers are more
// than cities, then part of the soldiers will not be distributed.
func (ac *AlienCommander) DistributeAliens() {
	var i int
	for name, city := range ac.WorldMap {
		if i >= len(ac.Soldiers) {
			break
		}
		city.Alien = ac.Soldiers[i]
		ac.WorldMap[name] = city
		i++
	}
	ac.Soldiers = ac.Soldiers[:i]
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

		if ac.iterations >= 10000 {
			fmt.Println("Stop the invasion because the 10000 iterations ware made.")
			return
		}

		availableAliens := len(ac.Soldiers) - ac.KilledSoldiers - int(ac.TrappedAliens.Load())
		if availableAliens < 2 {
			fmt.Printf("There is %d soldier left. Stop the invasion!\n", availableAliens)
			fmt.Println("Number of iterations: ", ac.iterations)
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
	FromAliens []Alien

	// CityName is the name of the city from where this report is send
	CityName string
}
