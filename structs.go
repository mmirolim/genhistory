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
func NewProfile(epochDate, regDate time.Time, edom string) Profile {
	// incr last user id
	id := atomic.AddInt64(&lid, 1)
	// email format id@provider.com
	mail := strconv.Itoa(int(id)) + "@" + edom + ".com"
	return Profile{Pid: id, RegDate: regDate, Email: mail, Edg: edom, Status: StatusAcitve}
}

type Population struct {
	m map[int64]Profile
	sync.Mutex
}

// add profile to population
func (pop *Population) Add(p Profile) {
	pop.Lock()
	pop.m[p.Pid] = p
	pop.Unlock()
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
func (p *Profile) Act(currentDate time.Time) []Event {
	var evs []Event
	if p.Pid == 0 || p.Status != 1 {
		return evs
	}
	// bernoulli distribution for unsub event
	unsubDist := rng.NewBernoulliGenerator(time.Now().UnixNano())
	// check if unsubscibed with unsub chance
	if unsubDist.Bernoulli_P(ur) {
		// update profile status on unsub
		p.Status = StatusUnsub
		p.UnsubDate = currentDate
		// gen one event for unsub
		evs = append(evs, Event{*p, 0.0, currentDate, 2})
		return evs
	}
	// he clicks n times
	clickDist := rng.NewWeibullGenerator(time.Now().UnixNano())
	nclicks := int(clickDist.Weibull(1, 1.5) * 3)
	costDist := rng.NewBernoulliGenerator(time.Now().UnixNano())
	// default click type is content
	clickCost := 0.2
	for i := 0; i < nclicks; i++ {
		// what cost of click
		if costDist.Bernoulli_P(advcr) {
			clickCost = 1.0
		}
		evs = append(evs, Event{*p, clickCost, currentDate, 1})
	}
	return evs
}

// NewPopulation generates N size population and
// saves it to db
func NewPopulation(n int, startDate time.Time) (*Population, error) {
	pop := &Population{}
	pop.m = make(map[int64]Profile)
	var p Profile
	var edom string
	var err error
	// domains will have uniform distribution
	uniProb := rng.NewUniformGenerator(time.Now().UnixNano())
	regDate := startDate
	for i := 0; i < n; i++ {
		// set random email domains
		edom = edg[uniProb.Int64n(int64(len(edg)))]
		// distribute uniformly registration date of profiles in 10 days before
		// start date
		regDate = startDate.Add(-time.Duration(uniProb.Int64n(10)) * 24 * time.Hour)
		p = NewProfile(startDate, regDate, edom)
		// add to population
		pop.Add(p)
		err = Save(&p)
		if err != nil {
			return pop, err
		}
	}
	return pop, nil
}
