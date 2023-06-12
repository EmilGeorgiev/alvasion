package app_test

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/EmilGeorgiev/alvasion/app"
	"github.com/stretchr/testify/assert"
)

type connection struct {
	DirectionName string // 0 - north, 1 - south, 2 - east, 3 -west
	City          string
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
			OutgoingRoadsNames: []string{"", "south=C5", "", "east=C1"},
		},
		{
			ID:                 3,
			Name:               "C3",
			OutgoingRoads:      []chan app.Alien{roads["c3c1"], roads["c3c6"], roads["c3c4"], nil},
			IncomingRoads:      []chan app.Alien{roads["c1c3"], roads["c6c3"], roads["c4c3"], nil},
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

// createAliens create a 9 aliens soldiers.
func createAliens() ([]app.Alien, chan app.Sitrep) {
	sitrep := make(chan app.Sitrep)
	return []app.Alien{
		{ID: 0},
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
		{ID: 7},
		{ID: 8},
	}, sitrep
}

type MockRandomizer struct {
	mock.Mock
}

func (m *MockRandomizer) ChooseRoad(roads []chan app.Alien) chan app.Alien {
	args := m.Called(roads)
	return args.Get(0).(chan app.Alien)
}

func mockCommandsOfCommander(m *MockRandomizer, roads map[string]chan app.Alien) {
	// In the beginning the map with aliens looks like that:
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C0,a0 |      | C1,a1 |      | C2,a2 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	 ↑  |        ↑  |        ↑  |
	//	 |  ↓        |  ↓        |  ↓
	//	|-------| ←--- |-------| ←--- |-------|
	//	| C3,a3 |      | C4,a4 |      | C5,a5 |
	//	|-------| ---→ |-------| ---→ |-------|
	//	 ↑  |        ↑  |        ↑  |
	//	 |  ↓        |  ↓        |  ↓
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

func TestStartInvasionWith6SoldiersAnd9Cities2(t *testing.T) {
	roads := createRoads()
	worldMap := createWorldMap(roads)
	aliens, sitrep := createAliens()
	buf := bytes.NewBufferString("")

	mockRand := new(MockRandomizer)
	mockCommandsOfCommander(mockRand, roads)

	commander := app.NewAlienCommander(worldMap, aliens, sitrep, mockRand, buf)
	commander.StartInvasion()
	actualReport := commander.GenerateReportForInvasion()

	expectedReport := "" +
		"C0 south=C3\n" +
		"C2 south=C5\n" +
		"C3 north=C0 south=C6\n" +
		"C5 north=C2\n" +
		"C6 north=C3\n"

	//line1, _ := buf.ReadString(0x0A)
	//line2, _ := buf.ReadString(0x0A)
	//line3, _ := buf.ReadString(0x0A)
	//line4, _ := buf.ReadString(0x0A)

	assert.Equal(t, expectedReport, actualReport)
	//assert.Equal(t, "C1 is destroyed from alien a0 and alien a2!", line1)
	//assert.Equal(t, "C4 is destroyed from alien a1 and alien a3!", line2)
	//assert.Equal(t, "C7 is destroyed from alien a6 and alien a8!", line3)
	//assert.Equal(t, "C8 is destroyed from alien a5 and alien a7!", line4)
}

//func TestStartInvasionWith6SoldiersAnd9Cities(t *testing.T) {
//	// SETUP
//	parts := make(chan []string)
//	go func() {
//		parts <- []string{"X1", "east=X2", "south=X4"}
//		parts <- []string{"X2", "east=X3", "west=X1", "south=X5"}
//		parts <- []string{"X3", "west=X2", "south=X6"}
//		parts <- []string{"X4", "east=X5", "north=X1", "south=X7"}
//		parts <- []string{"X5", "west=X4", "east=X6", "north=X2", "south=X8"}
//		parts <- []string{"X6", "west=X5", "north=X3", "south=X9"}
//		parts <- []string{"X7", "east=X8", "north=X4"}
//		parts <- []string{"X8", "west=X7", "east=X9", "north=X5"}
//		parts <- []string{"X9", "west=X8", "north=X6"}
//		close(parts)
//	}()
//
//	cities := buildCitiesAndTheirConnections()
//	wm := app.GenerateWorldMap(parts)
//	sitreps := make(chan app.Sitrep, 1)
//	aliens := []*app.Alien{
//		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
//		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps},
//	}
//	ac := app.NewAlienCommander(wm, aliens, sitreps)
//	notify := make(chan string)
//	done := make(chan struct{})
//	listenForNotificationsFromAlienCommander(cities, notify, done)
//
//	// ACTION
//	ac.SetNotifyDestroy(notify)
//	ac.StartInvasion()
//	close(notify)
//	<-done
//
//	// ASSERTIONS
//	actual := ac.GenerateReportForInvasion2()
//	expected := generateReport(cities)
//	assert.Equal(t, expected, actual)
//}
//
//func TestStartInvasionWith5AliensSoldiersAnd9Cities(t *testing.T) {
//	// SETUP
//	parts := make(chan []string)
//	go func() {
//		parts <- []string{"X1", "east=X2", "south=X4"}
//		parts <- []string{"X2", "east=X3", "west=X1", "south=X5"}
//		parts <- []string{"X3", "west=X2", "south=X6"}
//		parts <- []string{"X4", "east=X5", "north=X1", "south=X7"}
//		parts <- []string{"X5", "west=X4", "east=X6", "north=X2", "south=X8"}
//		parts <- []string{"X6", "west=X5", "north=X3", "south=X9"}
//		parts <- []string{"X7", "east=X8", "north=X4"}
//		parts <- []string{"X8", "west=X7", "east=X9", "north=X5"}
//		parts <- []string{"X9", "west=X8", "north=X6"}
//		close(parts)
//	}()
//
//	cities := buildCitiesAndTheirConnections()
//	wm := app.GenerateWorldMap(parts)
//	sitreps := make(chan app.Sitrep, 1)
//	aliens := []*app.Alien{
//		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps},
//		{ID: 3, Sitreps: sitreps}, {ID: 4, Sitreps: sitreps},
//	}
//	ac := app.NewAlienCommander(wm, aliens, sitreps)
//	notify := make(chan string)
//	done := make(chan struct{})
//	listenForNotificationsFromAlienCommander(cities, notify, done)
//
//	// ACTION
//	ac.SetNotifyDestroy(notify)
//	ac.StartInvasion()
//	close(notify)
//	<-done
//
//	// ASSERTIONS
//	actual := ac.GenerateReportForInvasion2()
//	expected := generateReport(cities)
//	assert.Equal(t, expected, actual)
//}
//
//func TestStartInvasionWith1AliensSoldiersAnd9Cities(t *testing.T) {
//	// SETUP
//	parts := make(chan []string)
//	go func() {
//		parts <- []string{"X1", "east=X2", "south=X4"}
//		parts <- []string{"X2", "east=X3", "west=X1", "south=X5"}
//		parts <- []string{"X3", "west=X2", "south=X6"}
//		parts <- []string{"X4", "east=X5", "north=X1", "south=X7"}
//		parts <- []string{"X5", "west=X4", "east=X6", "north=X2", "south=X8"}
//		parts <- []string{"X6", "west=X5", "north=X3", "south=X9"}
//		parts <- []string{"X7", "east=X8", "north=X4"}
//		parts <- []string{"X8", "west=X7", "east=X9", "north=X5"}
//		parts <- []string{"X9", "west=X8", "north=X6"}
//		close(parts)
//	}()
//
//	cities := buildCitiesAndTheirConnections()
//	wm := app.GenerateWorldMap(parts)
//	sitreps := make(chan app.Sitrep, 1)
//	aliens := []*app.Alien{{ID: 0, Sitreps: sitreps}}
//	ac := app.NewAlienCommander(wm, aliens, sitreps)
//	notify := make(chan string)
//	done := make(chan struct{})
//	listenForNotificationsFromAlienCommander(cities, notify, done)
//
//	// ACTION
//	ac.SetNotifyDestroy(notify)
//	ac.StartInvasion()
//	close(notify)
//	<-done
//
//	// ASSERTIONS
//	actual := ac.GenerateReportForInvasion2()
//	expected := generateReport(cities)
//	assert.Equal(t, expected, actual)
//}
//
//func TestStartInvasionWith11AliensSoldiersAnd9Cities(t *testing.T) {
//	parts := make(chan []string)
//	go func() {
//		parts <- []string{"X1", "east=X2", "south=X4"}
//		parts <- []string{"X2", "east=X3", "west=X1", "south=X5"}
//		parts <- []string{"X3", "west=X2", "south=X6"}
//		parts <- []string{"X4", "east=X5", "north=X1", "south=X7"}
//		parts <- []string{"X5", "west=X4", "east=X6", "north=X2", "south=X8"}
//		parts <- []string{"X6", "west=X5", "north=X3", "south=X9"}
//		parts <- []string{"X7", "east=X8", "north=X4"}
//		parts <- []string{"X8", "west=X7", "east=X9", "north=X5"}
//		parts <- []string{"X9", "west=X8", "north=X6"}
//		close(parts)
//	}()
//
//	cities := buildCitiesAndTheirConnections()
//	wm := app.GenerateWorldMap(parts)
//	sitreps := make(chan app.Sitrep, 1)
//	aliens := []*app.Alien{
//		{ID: 0, Sitreps: sitreps}, {ID: 1, Sitreps: sitreps}, {ID: 2, Sitreps: sitreps}, {ID: 3, Sitreps: sitreps},
//		{ID: 4, Sitreps: sitreps}, {ID: 5, Sitreps: sitreps}, {ID: 6, Sitreps: sitreps}, {ID: 7, Sitreps: sitreps},
//		{ID: 8, Sitreps: sitreps}, {ID: 9, Sitreps: sitreps}, {ID: 10, Sitreps: sitreps},
//	}
//	ac := app.NewAlienCommander(wm, aliens, sitreps)
//	notify := make(chan string)
//	done := make(chan struct{})
//	listenForNotificationsFromAlienCommander(cities, notify, done)
//
//	// ACTION
//	ac.SetNotifyDestroy(notify)
//	ac.StartInvasion()
//	close(notify)
//	<-done
//
//	// ASSERTIONS
//	actual := ac.GenerateReportForInvasion2()
//	expected := generateReport(cities)
//	assert.Equal(t, expected, actual)
//}
//
//func buildCitiesAndTheirConnections() map[string][]*connection {
//	return map[string][]*connection{
//		"X1": {{DirectionName: "south=X4", City: "X4"}, {DirectionName: "east=X2", City: "X2"}},
//		"X2": {{DirectionName: "south=X5", City: "X5"}, {DirectionName: "east=X3", City: "X3"}, {DirectionName: "west=X1", City: "X1"}},
//		"X3": {{DirectionName: "south=X6", City: "X6"}, {DirectionName: "west=X2", City: "X2"}},
//		"X4": {{DirectionName: "north=X1", City: "X1"}, {DirectionName: "south=X7", City: "X7"}, {DirectionName: "east=X5", City: "X5"}},
//		"X5": {{DirectionName: "north=X2", City: "X2"}, {DirectionName: "south=X8", City: "X8"}, {DirectionName: "east=X6", City: "X6"}, {DirectionName: "west=X4", City: "X4"}},
//		"X6": {{DirectionName: "north=X3", City: "X3"}, {DirectionName: "south=X9", City: "X9"}, {DirectionName: "west=X5", City: "X5"}},
//		"X7": {{DirectionName: "north=X4", City: "X4"}, {DirectionName: "east=X8", City: "X8"}},
//		"X8": {{DirectionName: "north=X5", City: "X5"}, {DirectionName: "east=X9", City: "X9"}, {DirectionName: "west=X7", City: "X7"}},
//		"X9": {{DirectionName: "north=X6", City: "X6"}, {DirectionName: "west=X8", City: "X8"}},
//	}
//}
//
//func listenForNotificationsFromAlienCommander(cities map[string][]*connection, notify chan string, done chan struct{}) {
//	go func() {
//		for n := range notify {
//			for _, conn := range cities[n] {
//				if conn == nil {
//					continue
//				}
//				conn2 := cities[conn.City]
//				for i, c := range conn2 {
//					if c == nil {
//						continue
//					}
//					if c.City == n {
//						conn2[i] = nil
//					}
//				}
//			}
//			delete(cities, n)
//		}
//		done <- struct{}{}
//	}()
//}
//
//func generateReport(cities map[string][]*connection) map[string][]string {
//	var keys []string
//	for k, _ := range cities {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//
//	result := map[string][]string{}
//	for _, name := range keys {
//		connections := cities[name]
//		for _, road := range connections {
//			if road == nil {
//				continue
//			}
//			result[name] = append(result[name], road.DirectionName)
//		}
//	}
//
//	return result
//}
