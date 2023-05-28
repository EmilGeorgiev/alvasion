package alvasion_test

import (
	"sync"
	"testing"

	"github.com/EmilGeorgiev/alvasion"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateCityDestructionWhenZeroAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{Name: "Foo"}
	wg := sync.WaitGroup{}

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	wg.Wait()

	// ASSERTIONS
	expected := alvasion.City{Name: "Foo"}
	assert.Equal(t, expected, c)
}

func TestEvaluateCityDestructionWhenOneAlienVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name:  "Foo",
		Roads: []chan alvasion.Alien{make(chan alvasion.Alien, 1)},
	}
	wg := sync.WaitGroup{}
	a := alvasion.Alien{Name: "55", Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[0] <- a

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actual := <-a.Sitreps
	wg.Wait()

	// ASSERTION
	expected := alvasion.Sitrep{
		From:            "55",
		CityName:        "Foo",
		IsCityDestroyed: false,
	}
	assert.Equal(t, expected, actual)
	assert.Equal(t, 1, len(c.Roads))
}

func TestEvaluateCityDestructionWhenTwoAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	a55 := alvasion.Alien{Name: "55", Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[0] <- a55
	a100 := alvasion.Alien{Name: "100", Sitreps: make(chan alvasion.Sitrep, 1)}
	c.Roads[1] <- a100

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep55 := <-a55.Sitreps
	actualRep100 := <-a100.Sitreps
	wg.Wait()

	// ASSERTION
	expectedRep55 := alvasion.Sitrep{From: "55", CityName: "Baz", IsCityDestroyed: true}
	expectedRep100 := alvasion.Sitrep{From: "100", CityName: "Baz", IsCityDestroyed: true}
	assert.Equal(t, expectedRep55, actualRep55)
	assert.Equal(t, expectedRep100, actualRep100)
}

func TestEvaluateCityDestructionWhenFourAliensVisitTheCity(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
			make(chan alvasion.Alien, 1),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 1)
	a1 := alvasion.Alien{Name: "1", Sitreps: reports}
	c.Roads[0] <- a1
	a2 := alvasion.Alien{Name: "2", Sitreps: reports}
	c.Roads[1] <- a2
	a3 := alvasion.Alien{Name: "3", Sitreps: reports}
	c.Roads[2] <- a3
	a4 := alvasion.Alien{Name: "4", Sitreps: reports}
	c.Roads[3] <- a4

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep1 := <-reports
	actualRep2 := <-reports
	actualRep3 := <-reports
	actualRep4 := <-reports
	wg.Wait()

	// ASSERTION
	expectedRep1 := alvasion.Sitrep{From: "1", CityName: "Baz", IsCityDestroyed: true}
	expectedRep2 := alvasion.Sitrep{From: "2", CityName: "Baz", IsCityDestroyed: true}
	expectedRep3 := alvasion.Sitrep{From: "3", CityName: "Baz", IsCityDestroyed: true}
	expectedRep4 := alvasion.Sitrep{From: "4", CityName: "Baz", IsCityDestroyed: true}
	assert.Equal(t, expectedRep1, actualRep1)
	assert.Equal(t, expectedRep2, actualRep2)
	assert.Equal(t, expectedRep3, actualRep3)
	assert.Equal(t, expectedRep4, actualRep4)
}

func TestEvaluateCityDestructionWhenTwoAliensVisitTheCityThroughTheSameChannel(t *testing.T) {
	// SETUP
	c := alvasion.City{
		Name: "Baz",
		Roads: []chan alvasion.Alien{
			make(chan alvasion.Alien, 2),
		},
	}
	wg := sync.WaitGroup{}
	reports := make(chan alvasion.Sitrep, 2)
	a1 := alvasion.Alien{Name: "1", Sitreps: reports}
	c.Roads[0] <- a1
	// Here the Sitrep channel is nil. This means that the test will panic if the Alien 2
	// push his report to the channel.
	a2 := alvasion.Alien{Name: "2", Sitreps: nil}
	c.Roads[0] <- a2

	// ACTION
	alvasion.EvaluateCityDestruction(&c, &wg)
	actualRep1 := <-reports
	wg.Wait()

	// ASSERTION
	expectedRep1 := alvasion.Sitrep{From: "1", CityName: "Baz", IsCityDestroyed: false}
	assert.Equal(t, expectedRep1, actualRep1)
}
