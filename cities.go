package alvasion

type City struct {
	Name  string
	South chan Alien
	West  chan Alien
	North chan Alien
	East  chan Alien
}
