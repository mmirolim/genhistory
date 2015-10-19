package main

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/leesper/go_rng"
)

var (
	StatusAcitve = 1
	StatusUnsub  = 2
)

type Profile struct {
	Pid     int64     // auto incremented counter
	Email   string    // id@edg.com
	Edg     string    // email provider
	RegDate time.Time // registration date 2006 Jan 02 resolution is day
	Status  int
}

// NewProfile generates Profiles according to predefined
// logic
func NewProfile(lifeStarted bool) Profile {
	var reg time.Time
	// incr last user id
	id := atomic.AddInt64(&lid, 1)
	// if pre live population set same date
	if !lifeStarted {
		reg = epS.Add(-10 * 24 * time.Hour)
	} else {
		// if life started
		reg = currentDate
	}
	// generate emails
	// domains will have uniform distribution
	uniProb := rng.NewUniformGenerator(time.Now().UnixNano())
	dom := edg[uniProb.Int64n(int64(len(edg)))]
	// email format id@provider.com
	mail := strconv.Itoa(int(id)) + "@" + dom + ".com"

	return Profile{Pid: id, RegDate: reg, Email: mail, Edg: dom}
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
	return Save(p)
}

// RM remvoe profile from population
func (pop *Population) RM(p Profile) (err error) {
	pop.Lock()
	err := Update(p, bson.M{"pid", p.Pid}, bson.M{})
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

// NewEvent create some random event for random profile
func NewEvent() Event {
	var ev Event
	return ev
}

// NewPopulation generates N size population and
// saves it to db
func NewPopulation(n int) {

}
