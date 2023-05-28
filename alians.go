package alvasion

import (
	"math/rand"
	"sync"
)

type AlienCommander struct {
	WorldMap                   WorldMap
	Soldiers                   []Alien
	NumberGettingTrappedAliens int
	Sitreps                    chan Sitrep

	wg sync.WaitGroup
}

func (ac *AlienCommander) GiveOrdersToTheAlienIn(c City) {
	// if in the city there is an alien, then the commander will give orders to him
	if c.Alien.Name == "" {
		return
	}

	// The commander selects a random active outgoing road and orders the alien to take that road.
	i := rand.Intn(len(c.Roads))
	c.Roads[i] <- c.Alien
}

func (ac *AlienCommander) StartNextIteration() {
	wg := sync.WaitGroup{}
	for _, city := range ac.WorldMap.Cities {
		ac.GiveOrdersToTheAlienIn(city)

		EvaluateCityDestruction(&city, &wg)
	}

	wg.Wait() // waiting the iteration of the invasion to finish

	// send signal to notify that the iteration is finished and the commander should prepare the next iteration

}

func (ac *AlienCommander) ListenForSitrep() {
	for sitrep := range ac.Sitreps {
		if sitrep.IsCityDestroyed {
			//ac.WorldMap.Cities[]
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
	ac.DistributeAliens()

	ac.StartNextIteration()
}

type Alien struct {
	Name string

	// Sitrep (Situation Report) is used to describe the current status of an ongoing mission.
	Sitreps chan Sitrep
}

type Sitrep struct {
	From            string
	CityName        string
	IsCityDestroyed bool
}
