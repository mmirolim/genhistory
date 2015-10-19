package main

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/leesper/go_rng"
)

var (
	StatusAcitve = 1
	StatusUnsub  = 2
)

type Profile struct {
	Pid       int64     // auto incremented counter
	Email     string    // id@edg.com
	Edg       string    // email provider
	RegDate   time.Time // registration date 2006 Jan 02 resolution is day
	UnsubDate time.Time // when profile unsubscribed
	Status    int
}

// Collection return mongo collection name
func (p *Profile) Collection() string {
	return "profiles"
}

// Index return index definition for mongo db
func (p *Profile) Index() []mgo.Index {
	return []mgo.Index{
		mgo.Index{
			Key: []string{"pid"},
		},
		mgo.Index{
			Key: []string{"edg"},
		},
		mgo.Index{
			Key: []string{"regdate"},
		},
		mgo.Index{
			Key: []string{"status"},
		},
	}
}

// NewProfile generates Profiles according to predefined
// logic
func NewProfile(epochDate time.Time, currentDay int, edom string) Profile {
	var reg time.Time
	// incr last user id
	id := atomic.AddInt64(&lid, 1)
	// if pre live population set same date
	if currentDay == 0 {
		reg = epochDate.Add(-10 * 24 * time.Hour)
	} else {
		// if life started
		reg = currentDate
	}
	// email format id@provider.com
	mail := strconv.Itoa(int(id)) + "@" + edom + ".com"
	return Profile{Pid: id, RegDate: reg, Email: mail, Edg: edom}
}

type Population struct {
	m map[int64]Profile
	sync.Mutex
}

// add profile to population
func (pop *Population) Add(p Profile) error {
	pop.Lock()
	pop.m[p.Pid] = p
	pop.Unlock()
	return Save(&p)
}

// RM remvoe profile from population
func (pop *Population) RM(p Profile) (err error) {
	pop.Lock()
	err = UpdateByPid(&p, p.Pid)
	if err != nil {
		return
	}
	delete(pop.m, p.Pid)
	pop.Unlock()
	return
}

// Get profile from population
func (pop *Population) Get(pid int64) Profile {
	return pop.m[pid]
}

// Size of population
func (pop *Population) Size() int {
	return len(pop.m)
}

type Event struct {
	Profile           // profile who generated event
	Cost    float64   // can be [1, 0.2, 0]
	Date    time.Time // resolution is day, event occured date
	Action  int       // what was an action [click, unsub]
}

// Collection return mongo collection name
func (ev *Event) Collection() string {
	return "events"
}

// Index return index definition for mongo db
func (ev *Event) Index() []mgo.Index {
	return []mgo.Index{
		mgo.Index{
			Key: []string{"date"},
		},
		mgo.Index{
			Key: []string{"action"},
		},
		mgo.Index{
			Key: []string{"profile.pid"},
		},
		mgo.Index{
			Key: []string{"profile.regdate"},
		},
		mgo.Index{
			Key: []string{"unsubdate"},
		},
	}
}

// Act generate slice of random events for a profile
func (p *Profile) Act() []Event {
	var evs []Event
	// bernoulli distribution for unsub event
	unsubDist := rng.NewBernoulliGenerator(time.Now().UnixNano())
	// check if unsubscibed with unsub chance
	if !unsubDist.Bernoulli_P(ur) {
		// gen one event for unsub
		evs = append(evs, Event{*p, 0.0, currentDate, 2})
		return evs
	}
	// he clicks n times
	clickDist := rng.NewWeibullGenerator(time.Now().UnixNano())
	nclicks := int(clickDist.Weibull(1, 1.5) * 3)
	costDist := rng.NewBernoulliGenerator(time.Now().UnixNano())
	clickCost := 1.0
	for i := 0; i < nclicks; i++ {
		// what cost of click
		if !costDist.Bernoulli_P(advcr) {
			clickCost = 0.2
		}
		evs = append(evs, Event{*p, clickCost, currentDate, 1})
	}
	return evs
}

// NewPopulation generates N size population and
// saves it to db
func NewPopulation(n int, startDate time.Time) (*Population, error) {
	var pop *Population
	var p Profile
	var edom string
	var err error
	// domains will have uniform distribution
	uniProb := rng.NewUniformGenerator(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		// generate emails
		edom = edg[uniProb.Int64n(int64(len(edg)))]
		p = NewProfile(startDate, 0, edom)
		// save to persistent storage
		err = Save(&p)
		if err != nil {
			return pop, err
		}
		// add to population
		pop.Add(p)
	}
	return pop, nil
}
