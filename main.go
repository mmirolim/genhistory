package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"xr/ltv/gendata/datastore"

	"github.com/leesper/go_rng"
)

var (
	// population number
	pN = 10000
	// epoch start
	epochStart = "2015-10-11"
	epS, _     = time.Parse("2006-01-02", epochStart)
	// current day for live a day func
	currentDate = epS
	// email domain group, profile email providers
	edg = []string{"yahoo", "aol", "hotmail", "google"}
	// registration rate
	rr = 0.005
	// unsubscribe rate
	ur = 0.002
	// click rate
	cr = 0.0338
	// there are two types of click weights
	adclick  = 1.0
	conclick = 0.2
	// number of history days to generate
	days = 42 // 6 weeks
	// last user id
	lid int64 = 0
)

type Profile struct {
	pid     int64     // auto incremented counter
	email   string    // id@edg.com
	edg     string    // email provider
	regDate time.Time // registration date 2006 Jan 02 resolution is day
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

	return Profile{pid: id, regDate: reg, email: mail, edg: dom}
}

type population struct {
	m map[int]Profile
	sync.Mutex
}

func init() {
	log.Println("connect to mongodb")
	err := ds.Init("localhost:27017", "click_history")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	fmt.Println("start generation")
	// on day 0 prepare population
	// generate profiles
	err = genProfiles(10)
	if err != nil {
		log.Fatalln("profiles generation failed", err)
	}
	// live a day from day first
	for d := 1; d < days; d++ {
		if err = liveADay(); err != nil {
			log.Println("liveADay encountered an error", err)
			break
		}

	}
}

// genProfiles creates random profiles and save collection
// to mongodb
func genProfiles(n int) error {
	var err error
	for i := 0; i < n; i++ {
		log.Printf("%+v\n", NewProfile(false))
	}
	return err
}

// one day from site audience
func liveADay() error {
	var err error

	return err
}
