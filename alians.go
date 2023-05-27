package alvasion

type AlienCommander struct {
}

type Alien struct {
	Name string

	// Sitrep (Situation Report) is used to describe the current status of an ongoing mission.
	Sitreps chan Sitrep
}

type Sitrep struct {
	CityName string
}