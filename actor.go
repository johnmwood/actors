package main

import (
	"fmt"
	"sync"
)

func NewAdserverActorSystem() *AdserverActorSystem {
	// initialize actors in slice of actors before creating actor system
	actors := []Actor{
		NewDecisionPublisher(),
	}
	return &AdserverActorSystem{
		actors: actors,
		tasks:  make(chan Task),
	}
}

func (sys *AdserverActorSystem) Run() {
	for _, actor := range sys.actors {
		go actor.Start()
	}

	for task := range sys.tasks {
		sys.SubmitTask(task)
	}
	sys.Shutdown()
}

func (sys *AdserverActorSystem) SubmitTask(task any) {
	switch task.(type) {
	case DecisionPublisherTask:
		sys.actors[DecisionActorType].AddTask(task)
	default:
		fmt.Errorf("task does not have associated actor: %+v", task)
	}
}

func (sys *AdserverActorSystem) Shutdown() {
	defer sys.wg.Done()
	for _, actor := range sys.actors {
		actor.Stop()
	}
	sys.wg.Wait()
	fmt.Printf("sys shutdown complete")
}

func NewDecisionPublisher() *DecisionPublisher {
	return &DecisionPublisher{
		id:       DecisionActorType,
		wg:       &sync.WaitGroup{},
		shutdown: make(chan bool),
		tasks:    make(chan DecisionPublisherTask),
	}
}

func (dp *DecisionPublisher) AddTask(task any) error {
	if t, ok := task.(DecisionPublisherTask); ok {
		dp.tasks <- t
		return nil
	}
	return fmt.Errorf("expected task type: DecisionPublisherTask. got: %t", task)
}

func (dp *DecisionPublisher) Start() {
	defer dp.wg.Done()
	dp.wg.Add(1)

	for task := range dp.tasks {
		task.Execute()
	}
	dp.shutdown <- true
}

func (dp *DecisionPublisher) Stop() {
	close(dp.tasks)
	<-dp.shutdown
}

func (task DecisionPublisherTask) Execute() {
	task.writeToPubSub()
}

func (task *DecisionPublisherTask) writeToPubSub() {
	fmt.Println("writing decision to Pub/Sub: %s", task.Decision)
}

func (d Decision) String() string {
	return fmt.Sprintf("decision kind %d - %s", d.kind, d.description)
}
