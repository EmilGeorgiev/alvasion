package alvasion

import "sync"

type AlienCommander struct {
	WorldMap                   WorldMap
	Soldiers                   []Alien
	NumberGettingTrappedAliens int
	Sitreps                    chan Sitrep

	wg sync.WaitGroup
}

func (ac *AlienCommander) GiveOrders() {

}

func (ac *AlienCommander) StartNextIteration() {

}

func (ac *AlienCommander) StopInvasion() {

}

func (ac *AlienCommander) StartInvasion() {
	for {

		//for i, city := range ac.WorldMap.Cities {
		//	city.Start
		//}
	}
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
