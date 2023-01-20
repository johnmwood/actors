package main

import (
	"math/rand"
	"sync"
	"time"
)

const (
	FirstCriteriaDecision = iota
	OtherDecision
)

/*
Order steps:
1. Order created
2. Order filled or rejected - i.e. decision made
  a. Maybe decision shouldn't be static? I.e. boolean or Do.Once().
	 Becase a order will normally have multiple decisions made before filling.
  b.
3. Decision sent to pub/sub
*/

func main() {
	decisions := []Decision{
		{
			kind:        FirstCriteriaDecision,
			description: "First criteria decision description",
		}, {
			kind:        OtherDecision,
			description: "Other criteria descriptions. This is different than the first criteria",
		}, {
			kind:        FirstCriteriaDecision,
			description: "another first criteria decision made",
		}, {
			kind:        OtherDecision,
			description: "another other decision",
		}, {
			kind:        FirstCriteriaDecision,
			description: "something",
		}, {
			kind:        OtherDecision,
			description: "something else",
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(decisions))

	actorSystem := NewAdserverActorSystem()
	go actorSystem.Run()

	for _, d := range decisions {
		SimulateMadeDecision(d, actorSystem, &wg)
	}
	wg.Wait()
}

func SimulateMadeDecision(d Decision, sys *AdserverActorSystem, wg *sync.WaitGroup) {
	task := DecisionPublisherTask{
		Decision: d,
	}
	go func(task Task, wg *sync.WaitGroup) {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		sys.SubmitTask(task)
		wg.Done()
	}(task, wg)
}
