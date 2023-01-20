package main

import "sync"

// these ids help determine the tasks to map to the correct actors
const (
	DecisionActorType = iota
	ImpressionActorType
)

type Actor interface {
	AddTask(task any) error
	Start()
	Stop()
}

type Task interface {
	Execute()
}

type Publisher interface {
	writeToPubSub()
}

type PublisherTask interface {
	Publisher
	Task
}

type ActorSystem interface {
	Run()
	SubmitTask(task any)
	Shutdown()
}

type AdserverActorSystem struct {
	actors []Actor
	tasks  chan Task
	wg     *sync.WaitGroup
}

type Decision struct {
	kind        int
	description string
}

type DecisionPublisher struct {
	id       int
	wg       *sync.WaitGroup
	shutdown chan bool
	tasks    chan DecisionPublisherTask
}

type DecisionPublisherTask struct {
	Decision
	/* pubsubCreds, pubsubVars */
}
