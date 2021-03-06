package main

import (
	"fmt"
	"log"
	"time"
)

var (
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
	// chance of clicking adv
	advcr = 0.70
	// number of history days to generate
	days = 20 // days to generate events
	// population size
	popSize = 100
	// last user id
	lid int64 = 0
	// Population
	Pop *Population
)

func init() {
	log.Println("connect to mongodb")
	err := mongoInit("localhost:27017", "ltv", &Profile{}, &Event{})
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	fmt.Println("start generation")
	// on day 0 prepare population
	// generate profiles
	Pop, err = NewPopulation(popSize, epS)
	log.Println("generated population", Pop)
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

// one day from site audience
func liveADay() error {
	var err error
	var events []Event
	var p Profile
	log.Println("Day number", currentDate)
	// range through population and let profiles to act
	for i := 1; i <= popSize; i++ {
		p = Pop.Get(int64(i))
		events = p.Act(currentDate)
		// update profile if status changed
		if p.Status == StatusUnsub {
			Pop.Add(p)
			err = UpdateByPid(&p, p.Pid)
			if err != nil {
				log.Println(err)
			}

		}
		// range over events and save them
		for _, ev := range events {
			err = Save(&ev)
			if err != nil {
				log.Println(err)
			}
		}

	}
	currentDate = currentDate.Add(24 * time.Hour)
	return err
}
