package alvasion

import (
	"fmt"
	"math/rand"
	"sync"
)

type AlienCommander struct {
	WorldMap       WorldMap
	Soldiers       []Alien
	TrappedAliens  int
	KilledSoldiers int
	Sitreps        chan Sitrep

	iterations    int
	iterationDone chan struct{}
	wg            sync.WaitGroup
}

func (ac *AlienCommander) GiveOrdersToTheAlienIn(c City) {
	// if in the city there is an alien, then the commander will give orders to him
	if c.Alien.Sitreps == nil {
		return
	}

	// The commander selects a random active outgoing road and orders the alien to take that road.
	i := rand.Intn(len(c.Roads))
	c.Roads[i] <- c.Alien
}

func (ac *AlienCommander) StartNextIteration() {
	wg := sync.WaitGroup{}
	// give orders to all soldiers
	for _, city := range ac.WorldMap.Cities {
		ac.GiveOrdersToTheAlienIn(*city)
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
	ac.iterationDone <- struct{}{} // this will stop the sitrep listener, because no more reports will be sent

	// when a city is destroyed all roads leading out or in of the town also should be destroyed. It is important
	// to keep in mind that one road always connect two different cities
	for _, city := range ac.WorldMap.Cities {
		EvaluateRoadsDestruction(city, &wg)
	}
	wg.Wait()
}

func (ac *AlienCommander) ListenForSitrep() {
	for {
		select {
		case report := <-ac.Sitreps:
			if !report.IsCityDestroyed {
				break
			}
			ac.WorldMap.CleanCity(report.CityName)
			ac.Soldiers[report.From].Killed = true
		default:
			// if there is no reports maybe the current iteration of
			// the invasion finished and all soldiers send their reports.
		}

		select {
		case <-ac.iterationDone:
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

func (ac *AlienCommander) DistributeAliens() {
	var i int
	for _, city := range ac.WorldMap.Cities {
		i++
		if i >= len(ac.Soldiers) {
			break
		}
		city.Alien = ac.Soldiers[i]
	}
}

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

		if (ac.KilledSoldiers + ac.TrappedAliens) >= len(ac.Soldiers)-1 {
			fmt.Println("Stop the invasion because less then 2 aliens are alive and not getting trapped.")
			return
		}
	}
}

type Alien struct {
	ID int

	// Sitrep (Situation Report) is used to describe the current status of an ongoing mission.
	Sitreps chan Sitrep

	Killed  bool
	Trapped bool
}

type Sitrep struct {
	From            int
	CityName        string
	IsCityDestroyed bool
}
