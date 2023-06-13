package app_test

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/EmilGeorgiev/alvasion/app"
	"github.com/stretchr/testify/assert"
)

func TestStartInvasionWith9SoldiersAnd9Cities(t *testing.T) {
	roads := createRoads()
	worldMap := createWorldMap(roads)
	aliens := []app.Alien{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}}
	buf := bytes.NewBufferString("")

	mockRand := new(MockRandomizer)
	mockMovementsOfThe9Aliens(mockRand, roads)

	commander := app.NewAlienCommander(worldMap, aliens, mockRand, buf, 10000)
	commander.StartInvasion()
	actualReport := commander.GenerateReportForInvasion()

	expectedReport := "" +
		"C0 south=C3\n" +
		"C2 south=C5\n" +
		"C3 north=C0 south=C6\n" +
		"C5 north=C2\n" +
		"C6 north=C3\n"

	assert.Equal(t, expectedReport, actualReport)
	assert.Contains(t, buf.String(), "C1 is destroyed from alien 0 and alien 2!")
	assert.Contains(t, buf.String(), "C4 is destroyed from alien 1 and alien 3!")
	assert.Contains(t, buf.String(), "C7 is destroyed from alien 6 and alien 8!")
	assert.Contains(t, buf.String(), "C8 is destroyed from alien 5 and alien 7!")
}

func TestStartInvasionWith6SoldiersAnd9Cities(t *testing.T) {
	roads := createRoads()
	worldMap := createWorldMap(roads)
	aliens := []app.Alien{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	buf := bytes.NewBufferString("")

	mockRand := new(MockRandomizer)
	mockMovementsOfThe6Aliens(mockRand, roads)

	commander := app.NewAlienCommander(worldMap, aliens, mockRand, buf, 10000)
	commander.StartInvasion()
	actualReport := commander.GenerateReportForInvasion()

	expectedReport := "" +
		"C0 south=C3 east=C1\n" +
		"C1 south=C4 east=C2 west=C0\n" +
		"C2 south=C5 west=C1\n" +
		"C3 north=C0 east=C4\n" +
		"C4 north=C1 east=C5 west=C3\n" +
		"C5 north=C2 south=C8 west=C4\n" +
		"C8 north=C5\n"

	assert.Equal(t, expectedReport, actualReport)
	assert.Contains(t, buf.String(), "C6 is destroyed from alien 0 and alien 4!")
	assert.Contains(t, buf.String(), "C7 is destroyed from alien 1 and alien 3 and alien 5!")
}

func TestStopInvasionBecauseReachedMaximumNumberOfIterations(t *testing.T) {
	roads := createRoads()
	worldMap := createWorldMap(roads)
	aliens := []app.Alien{{ID: 0}, {ID: 1}}
	buf := bytes.NewBufferString("")

	mockRand := new(MockRandomizer)
	mockMovementsOfThe2Aliens(mockRand, roads)

	maxNumberOfIterations := 2
	commander := app.NewAlienCommander(worldMap, aliens, mockRand, buf, maxNumberOfIterations)
	commander.StartInvasion()
	actualReport := commander.GenerateReportForInvasion()

	expectedReport := "" +
		"C0 south=C3 east=C1\n" +
		"C1 south=C4 east=C2 west=C0\n" +
		"C2 south=C5 west=C1\n" +
		"C3 north=C0 south=C6 east=C4\n" +
		"C4 north=C1 south=C7 east=C5 west=C3\n" +
		"C5 north=C2 south=C8 west=C4\n" +
		"C6 north=C3 east=C7\n" +
		"C7 north=C4 east=C8 west=C6\n" +
		"C8 north=C5 west=C7\n"

	assert.Equal(t, expectedReport, actualReport)
	assert.Equal(t, "", buf.String())
}

// createRoads create all roads between cities (incoming and outgoing). For the test we will use 9 cities
// Here is an example of cities and roads (C0, C1, ... C8 are the name of the cities):
//
//	|----| ←--- |----| ←--- |----|
//	| C0 |      | C1 |      | C2 |
//	|----| ---→ |----| ---→ |----|
//	 ↑  |        ↑  |        ↑  |
//	 |  ↓        |  ↓        |  ↓
//	|----| ←--- |----| ←--- |----|
//	| C3 |      | C4 |      | C5 |
//	|----| ---→ |----| ---→ |----|
//	 ↑  |        ↑  |        ↑  |
//	 |  ↓        |  ↓        |  ↓
//	|----| ←--- |----| ←--- |----|
//	| C6 |      | C7 |      | C8 |
//	|----| ---→ |----| ---→ |----|
func createRoads() map[string]chan app.Alien {
	return map[string]chan app.Alien{
		"c0c1": make(chan app.Alien, 1),
		"c0c3": make(chan app.Alien, 1),

		"c1c0": make(chan app.Alien, 1),
		"c1c2": make(chan app.Alien, 1),
		"c1c4": make(chan app.Alien, 1),

		"c2c1": make(chan app.Alien, 1),
		"c2c5": make(chan app.Alien, 1),
		"c3c0": make(chan app.Alien, 1),
		"c3c4": make(chan app.Alien, 1),
		"c3c6": make(chan app.Alien, 1),

		"c4c1": make(chan app.Alien, 1),
		"c4c3": make(chan app.Alien, 1),
		"c4c5": make(chan app.Alien, 1),
		"c4c7": make(chan app.Alien, 1),

		"c5c2": make(chan app.Alien, 1),
		"c5c4": make(chan app.Alien, 1),
		"c5c8": make(chan app.Alien, 1),

		"c6c3": make(chan app.Alien, 1),
		"c6c7": make(chan app.Alien, 1),

		"c7c4": make(chan app.Alien, 1),
		"c7c6": make(chan app.Alien, 1),
		"c7c8": make(chan app.Alien, 1),

		"c8c5": make(chan app.Alien, 1),
		"c8c7": make(chan app.Alien, 1),
	}
}

func createWorldMap(roads map[string]chan app.Alien) []app.City {
	return []app.City{
		{
			ID:                 0,
			Name:               "C0",
			OutgoingRoads:      []chan app.Alien{nil, roads["c0c3"], roads["c0c1"], nil},
			IncomingRoads:      []chan app.Alien{nil, roads["c3c0"], roads["c1c0"], nil},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"", "south=C3", "east=C1", ""},
		},
		{
			ID:                 1,
			Name:               "C1",
			OutgoingRoads:      []chan app.Alien{nil, roads["c1c4"], roads["c1c2"], roads["c1c0"]},
			IncomingRoads:      []chan app.Alien{nil, roads["c4c1"], roads["c2c1"], roads["c0c1"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"", "south=C4", "east=C2", "west=C0"},
		},
		{
			ID:                 2,
			Name:               "C2",
			OutgoingRoads:      []chan app.Alien{nil, roads["c2c5"], nil, roads["c2c1"]},
			IncomingRoads:      []chan app.Alien{nil, roads["c5c2"], nil, roads["c1c2"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"", "south=C5", "", "west=C1"},
		},
		{
			ID:                 3,
			Name:               "C3",
			OutgoingRoads:      []chan app.Alien{roads["c3c0"], roads["c3c6"], roads["c3c4"], nil},
			IncomingRoads:      []chan app.Alien{roads["c0c3"], roads["c6c3"], roads["c4c3"], nil},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C0", "south=C6", "east=C4", ""},
		},
		{
			ID:                 4,
			Name:               "C4",
			OutgoingRoads:      []chan app.Alien{roads["c4c1"], roads["c4c7"], roads["c4c5"], roads["c4c3"]},
			IncomingRoads:      []chan app.Alien{roads["c1c4"], roads["c7c4"], roads["c5c4"], roads["c3c4"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C1", "south=C7", "east=C5", "west=C3"},
		},
		{
			ID:                 5,
			Name:               "C5",
			OutgoingRoads:      []chan app.Alien{roads["c5c2"], roads["c5c8"], nil, roads["c5c4"]},
			IncomingRoads:      []chan app.Alien{roads["c2c5"], roads["c8c5"], nil, roads["c4c5"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C2", "south=C8", "", "west=C4"},
		},
		{
			ID:                 6,
			Name:               "C6",
			OutgoingRoads:      []chan app.Alien{roads["c6c3"], nil, roads["c6c7"], nil},
			IncomingRoads:      []chan app.Alien{roads["c3c6"], nil, roads["c7c6"], nil},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C3", "", "east=C7", ""},
		},
		{
			ID:                 7,
			Name:               "C7",
			OutgoingRoads:      []chan app.Alien{roads["c7c4"], nil, roads["c7c8"], roads["c7c6"]},
			IncomingRoads:      []chan app.Alien{roads["c4c7"], nil, roads["c8c7"], roads["c6c7"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C4", "", "east=C8", "west=C6"},
		},
		{
			ID:                 8,
			Name:               "C8",
			OutgoingRoads:      []chan app.Alien{roads["c8c5"], nil, nil, roads["c8c7"]},
			IncomingRoads:      []chan app.Alien{roads["c5c8"], nil, nil, roads["c7c8"]},
			IsDestroyed:        false,
			Alien:              nil,
			OutgoingRoadsNames: []string{"north=C5", "", "", "west=C7"},
		},
	}
}

func mockMovementsOfThe9Aliens(m *MockRandomizer, roads map[string]chan app.Alien) {
	// In the beginning the map with aliens looks like that:
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0,a0 |      | C1,a1 |      | C2,a2 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3,a3 |      | C4,a4 |      | C5,a5 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C6,a6 |      | C7,a7 |      | C8,a8 |
	//	|-------| ---→ |-------| ---→ |-------|
	//
	// every city has one alien. The the commander gives random orders to the soldiers:

	m.On("ChooseRoad", mock.Anything).Return(roads["c0c1"]).Once() // first alien (a0) move from C0 to C1
	m.On("ChooseRoad", mock.Anything).Return(roads["c1c4"]).Once() // second alien (a1) move from C1 to C4
	m.On("ChooseRoad", mock.Anything).Return(roads["c2c1"]).Once() // third alien (a2) move from C2 to C1
	m.On("ChooseRoad", mock.Anything).Return(roads["c3c4"]).Once() // fourth alien (a3) move from C3 to C4
	m.On("ChooseRoad", mock.Anything).Return(roads["c4c5"]).Once() // fifth alien (a4) move from C4 to C5
	m.On("ChooseRoad", mock.Anything).Return(roads["c5c8"]).Once() // sixth alien (a5) move from C5 to C8
	m.On("ChooseRoad", mock.Anything).Return(roads["c6c7"]).Once() // seventh alien (a6) move from C6 to C7
	m.On("ChooseRoad", mock.Anything).Return(roads["c7c8"]).Once() // eighth alien (a7) move from C7 to C8
	m.On("ChooseRoad", mock.Anything).Return(roads["c8c7"]).Once() // ninth alien (a8) move from C8 to C7

	// after above moving of aliens the map will look like that:
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C0    |      | C1,a0,a2 |      | C2       |
	//	|-------| ---→ |----------| ---→ |----------|
	//	 ↑  |           ↑        |         ↑    |
	//	 |  ↓           |        ↓         |    ↓
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C3    |      | C4,a1,a3 |      | C5,a4    |
	//	|-------| ---→ |----------| ---→ |----------|
	//	 ↑  |           ↑        |         ↑    |
	//	 |  ↓           |        ↓         |    ↓
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C6    |      | C7,a6,a8 |      | C8,a5,a7 |
	//	|-------| ---→ |----------| ---→ |----------|
	//
	// This means that the cities C1, C4, C7 and C8 will be destroyed because they have two aliens at the
	// same time and the map after the iteration will looks like that:
	//	|-------|                        |----------|
	//	| C0    |                        | C2       |
	//	|-------|                        |----------|
	//	 ↑  |                              ↑    |
	//	 |  ↓                              |    ↓
	//	|-------|                        |----------|
	//	| C3    |                        | C5,a4    |
	//	|-------|                        |----------|
	//	 ↑  |
	//	 |  ↓
	//	|-------|
	//	| C6    |
	//	|-------|
	//
	// after that the commander MUST stop the invasion because there is only one alien and he can't destroy a city alone.
}

func mockMovementsOfThe6Aliens(m *MockRandomizer, roads map[string]chan app.Alien) {
	// In the beginning the map with aliens looks like that:
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0,a0 |      | C1,a1 |      | C2,a2 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |            ↑  |           ↑  |
	//	  |  ↓            |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3,a3 |      | C4,a4 |      | C5,a5 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C6    |      | C7    |      | C8    |
	//	|-------| ---→ |-------| ---→ |-------|
	//
	// the first 6 cities have one alien. The the commander gives this orders to the soldiers:

	m.On("ChooseRoad", mock.Anything).Return(roads["c0c3"]).Once() // first alien (a0) move from C0 to C3
	m.On("ChooseRoad", mock.Anything).Return(roads["c1c4"]).Once() // second alien (a1) move from C1 to C4
	m.On("ChooseRoad", mock.Anything).Return(roads["c2c5"]).Once() // third alien (a2) move from C2 to C5
	m.On("ChooseRoad", mock.Anything).Return(roads["c3c6"]).Once() // fourth alien (a3) move from C3 to C6
	m.On("ChooseRoad", mock.Anything).Return(roads["c4c7"]).Once() // fifth alien (a4) move from C4 to C7
	m.On("ChooseRoad", mock.Anything).Return(roads["c5c8"]).Once() // sixth alien (a5) move from C5 to C8

	// after above moving of aliens the map will look like that:
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C0    |      | C1       |      | C2       |
	//	|-------| ---→ |----------| ---→ |----------|
	//	 ↑  |           ↑        |         ↑    |
	//	 |  ↓           |        ↓         |    ↓
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C3,a0 |      | C4,a1    |      | C5,a2    |
	//	|-------| ---→ |----------| ---→ |----------|
	//	 ↑  |           ↑        |         ↑    |
	//	 |  ↓           |        ↓         |    ↓
	//	|-------| ←--- |----------| ←--- |----------|
	//	| C6,a3 |      | C7,a4    |      | C8,a5    |
	//	|-------| ---→ |----------| ---→ |----------|
	//
	// This means that after the first iteration No cities will be destroyed because they all have one or zero aliens.

	// In the second iteration the commander gives these orders:
	m.On("ChooseRoad", mock.Anything).Return(roads["c3c6"]).Once() // first alien (a0) move from C3 to C6
	m.On("ChooseRoad", mock.Anything).Return(roads["c4c7"]).Once() // second alien (a1) move from C4 to C7
	m.On("ChooseRoad", mock.Anything).Return(roads["c5c8"]).Once() // third alien (a2) move from C5 to C8
	m.On("ChooseRoad", mock.Anything).Return(roads["c6c7"]).Once() // fourth alien (a3) move from C6 to C7
	m.On("ChooseRoad", mock.Anything).Return(roads["c7c6"]).Once() // fifth alien (a4) move from C7 to C6
	m.On("ChooseRoad", mock.Anything).Return(roads["c8c7"]).Once() // sixth alien (a5) move from C8 to C7

	// after that the aliens will be places like that:
	//	|----------| ←--- |-------------| ←--- |----------|
	//	| C0       |      | C1          |      | C2       |
	//	|----------| ---→ |-------------| ---→ |----------|
	//	   ↑  |             ↑    |               ↑    |
	//	   |  ↓             |    ↓               |    ↓
	//	|----------| ←--- |-------------| ←--- |----------|
	//	| C3       |      | C4          |      | C5       |
	//	|----------| ---→ |-------------| ---→ |----------|
	//	   ↑  |              ↑      |             ↑    |
	//	   |  ↓              |      ↓             |    ↓
	//	|----------| ←--- |-------------| ←--- |----------|
	//	| C6,a0,a4 |      | C7,a1,a3,a5 |      | C8,a2    |
	//	|----------| ---→ |-------------| ---→ |----------|
	//
	// this means that the cities C6 and C7 will be destroyed.

	//	|----------| ←--- |-------------| ←--- |----------|
	//	| C0       |      | C1          |      | C2       |
	//	|----------| ---→ |-------------| ---→ |----------|
	//	   ↑  |             ↑    |               ↑    |
	//	   |  ↓             |    ↓               |    ↓
	//	|----------| ←--- |-------------| ←--- |----------|
	//	| C3       |      | C4          |      | C5       |
	//	|----------| ---→ |-------------| ---→ |----------|
	//	                                          ↑    |
	//	                                          |    ↓
	//	                                       |----------|
	//	                                       | C8,a2    |
	//	                                       |----------|
	// after that the commander MUST stop the invasion because there is only one alien (a2) and he can't destroy a city alone.
}

func mockMovementsOfThe2Aliens(m *MockRandomizer, roads map[string]chan app.Alien) {
	// In the beginning the map with aliens looks like that:
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0,a0 |      | C1,a1 |      | C2    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |            ↑  |           ↑  |
	//	  |  ↓            |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3    |      | C4    |      | C5    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C6    |      | C7    |      | C8    |
	//	|-------| ---→ |-------| ---→ |-------|
	//
	// the first 2 cities have one alien. The the commander gives this orders to the soldiers:

	m.On("ChooseRoad", mock.Anything).Return(roads["c0c3"]).Once() // first alien (a0) move from C0 to C3
	m.On("ChooseRoad", mock.Anything).Return(roads["c1c4"]).Once() // second alien (a1) move from C1 to C4

	// after that the map will looks like that:
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0    |      | C1    |      | C2    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |            ↑  |           ↑  |
	//	  |  ↓            |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3,a0 |      | C4,a1 |      | C5    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C6    |      | C7    |      | C8    |
	//	|-------| ---→ |-------| ---→ |-------|
	//

	m.On("ChooseRoad", mock.Anything).Return(roads["c3c0"]).Once() // first alien (a0) move from C3 to C0
	m.On("ChooseRoad", mock.Anything).Return(roads["c4c1"]).Once() // second alien (a1) move from C4 to C1

	// now the map looks like that
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0,a0 |      | C1,a1 |      | C2    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |            ↑  |           ↑  |
	//	  |  ↓            |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3    |      | C4    |      | C5    |
	//	|-------| ---→ |-------| ---→ |-------|
	//	  ↑  |           ↑  |           ↑  |
	//	  |  ↓           |  ↓           |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C6    |      | C7    |      | C8    |
	//	|-------| ---→ |-------| ---→ |-------|
	//
}

type MockRandomizer struct {
	mock.Mock
}

func (m *MockRandomizer) ChooseRoad(roads []chan app.Alien) chan app.Alien {
	args := m.Called(roads)
	return args.Get(0).(chan app.Alien)
}
