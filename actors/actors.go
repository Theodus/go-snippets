package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	actor :=
		newActor("yup", 1,
			newNotify(func(actor *someActor) { done <- struct{}{} }))
	actor.messagec <- func() {
		fmt.Println(actor.state)
	}
	time.Sleep(time.Second)
	actor.quitc <- struct{}{}
	<-done
	fmt.Println("complete")
}

type actorNotify interface {
	shutdown(actor *someActor)
}

func newNotify(shutdown func(actor *someActor)) actorNotify {
	return dummyNotify{
		shutdownf: shutdown,
	}
}

type dummyNotify struct {
	shutdownf func(actor *someActor)
}

func (notify dummyNotify) shutdown(actor *someActor) {
	notify.shutdownf(actor)
}

type someActor struct {
	state    string
	notify   actorNotify
	messagec chan func()
	quitc    chan struct{}
}

func newActor(state string, buffer int, notify actorNotify) *someActor {
	actor := &someActor{
		state:    state,
		notify:   notify,
		messagec: make(chan func(), buffer),
		quitc:    make(chan struct{}),
	}
	go actor.loop()
	return actor
}

func (actor *someActor) loop() {
	for {
		select {
		case be := <-actor.messagec:
			be()
		case <-actor.quitc:
			actor.notify.shutdown(actor)
			return
		}
	}
}
