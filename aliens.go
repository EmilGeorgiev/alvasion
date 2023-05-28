package alvasion

import (
	"bytes"
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
	IterationDone chan struct{}
	wg            sync.WaitGroup
}

//var idxDirectionMap = map[int]string{
//	0: "north"
//}

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
			row = " " + road
		}
		buf.WriteString(row + "\n")
	}

	return ""
}

func NewAlienCommander(wm WorldMap, aliens []Alien, sitreps chan Sitrep) AlienCommander {
	return AlienCommander{
		WorldMap:       wm,
		Soldiers:       aliens,
		TrappedAliens:  0,
		KilledSoldiers: 0,
		Sitreps:        sitreps,
		IterationDone:  make(chan struct{}),
	}
}

func (ac *AlienCommander) GiveOrdersToTheAlienIn(c City) {
	// if in the city there is an alien, then the commander will give orders to him
	if c.Alien.Sitreps == nil {
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

	availableRoads[i] <- c.Alien

}

func (ac *AlienCommander) StartNextIteration() {
	wg := sync.WaitGroup{}
	fmt.Println("00000000000000000000000000")
	// give orders to all soldiers
	for _, city := range ac.WorldMap.Cities {
		ac.GiveOrdersToTheAlienIn(*city)
	}

	fmt.Println("1111111111111111111111111111")
	// prepare to listen for incoming situation reports about the evaluation of the invasion
	go ac.ListenForSitrep()

	fmt.Println("22222222222222222222222")
	// After issuing all orders, the commander can evaluate the consequences and assess the
	// damage inflicted upon the cities as a result of these commands.
	for _, city := range ac.WorldMap.Cities {
		EvaluateCityDestruction(city, &wg)
	}

	fmt.Println("33333333333333333333")
	wg.Wait() // waiting the current iteration of the invasion to finish
	fmt.Println("4444444444444444444")
	// send signal to notify that the iteration is finished. The commander should prepare the next iteration.
	ac.IterationDone <- struct{}{} // this will stop the sitrep listener, because no more reports will be sent
	fmt.Println("55555555555555555")
	// when a city is destroyed all roads leading out or in of the town also should be destroyed. It is important
	// to keep in mind that one road always connect two different cities
	for _, city := range ac.WorldMap.Cities {
		EvaluateRoadsDestruction(city, &wg)
	}
	fmt.Println("666666666666666666666666")
	wg.Wait()
	fmt.Println("777777777777777777")

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
