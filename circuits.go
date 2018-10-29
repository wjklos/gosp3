package main

import (
	"fmt"
	"net/http"
)

const (
	recvSlowPct = .75
	recvStopPct = .95
	holdSlowPct = .50
	holdStopPct = .95
)

var ()

// Circuit ...
type Circuit struct {
	Conductor Component
	Switch    Component
	Send      *Component
}

// CircuitInterface ...
type CircuitInterface interface {
	New()
}

// Component ...
type Component struct {
	Name           string
	Channel        chan string
	Depth          int
	Notifier       chan bool
	LastOperation  int
	SlowPercentage float64
	StopPercentage float64
}

// ComponentInterface ...
type ComponentInterface interface {
	New()
	Check()
	Notify()
}

func init() {
}

// New ...
func (c *Circuit) New() *Circuit {
	c.Conductor = *c.Conductor.New("receiver")
	c.Switch = *c.Switch.New("hold")
	c.Send = nil
	return c
}

// New ...
func (cc *Component) New(name string) *Component {
	cc.Name = name
	cc.Depth = maxMessages
	cc.Channel = make(chan string, cc.Depth)
	cc.Notifier = make(chan bool, 1)
	cc.LastOperation = http.StatusAccepted
	cc.SlowPercentage = holdSlowPct
	cc.StopPercentage = holdStopPct
	return cc
}

// Check ...
func (cc *Component) Check() chan bool {
	return cc.Notifier
}

// Notify ...
func (cc *Component) Notify() {
	// Since this is only a "call to action" channel, it only needs one call.
	// If there is already a message in it, then someone else made that call.
	if len(cc.Channel) < cap(cc.Channel) {
		cc.Notifier <- true
	}
}

// Fill ...
func (cc *Component) Fill(item string) bool {
	// Make sure there is room in the list before adding any thing to it.
	if (float64(len(cc.Channel)) / float64(cap(cc.Channel))) < cc.StopPercentage {
		cc.Channel <- item
		cc.Notify()
		fmt.Printf("%s list: (%s) %5.2f <? %5.2f\n", cc.Name, item, (float64(len(cc.Channel)) / float64(cap(cc.Channel))), cc.StopPercentage)
		if (float64(len(cc.Channel)) / float64(cap(cc.Channel))) < cc.SlowPercentage {
			cc.LastOperation = http.StatusAccepted
			return true
		} // if
		cc.LastOperation = http.StatusTooManyRequests
		return true
	} // if
	cc.LastOperation = http.StatusServiceUnavailable
	return false
} // func

// Deplete ...
func (cc *Component) Deplete() string {
	return <-cc.Channel
}
