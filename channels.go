package main

import (
	"fmt"
	"net/http"
)

const (
	recvSlowPct = .75
	recvStopPct = .95
	holdSlowPct = .75
	holdStopPct = .95
)

var (
	// ReceiveNotifier ...
	ReceiveNotifier chan bool
	// SendNotifier ...
	SendNotifier chan bool
	// HoldNotifier ...
	HoldNotifier chan bool
	// ReceiveList ...
	ReceiveList chan string
	// SendList ...
	SendList chan string
	// HoldList ...
	HoldList chan string
	// ReceiveLastOp ...
	ReceiveLastOp chan int
	// SendLastOp ...
	SendLastOp chan int
)

func init() {
	ReceiveNotifier = make(chan bool, 1)
	ReceiveList = make(chan string, maxMessages)
	ReceiveLastOp = make(chan int, 1)

	SendList = make(chan string, maxMessages)
	SendLastOp = make(chan int, 1)

	HoldList = make(chan string, maxMessages)
}

// GetReceiveNotifier ...
func GetReceiveNotifier() chan bool {
	return ReceiveNotifier
}

// NotifyReceiveList ...
func NotifyReceiveList() {
	// Since this is only a "call to action" channel, it only needs one call.
	// If there is already a message in it, then someone else made that call.
	if len(ReceiveNotifier) < cap(ReceiveNotifier) {
		ReceiveNotifier <- true
	}
}

// FillReceiveList ...
func FillReceiveList(item string) bool {
	// Make sure there is room in the list before adding any thing to it.
	if (float64(len(ReceiveList)) / float64(cap(ReceiveList))) < float64(recvStopPct) {
		fmt.Printf("ITEM: %s\n", item)
		ReceiveList <- item
		NotifyReceiveList()
		fmt.Printf("RECEIVE LIST: %5.2f <? %5.2f\n", (float64(len(ReceiveList)) / float64(cap(ReceiveList))), float64(recvStopPct))
		if (float64(len(ReceiveList)) / float64(cap(ReceiveList))) < float64(recvSlowPct) {
			ReceiveLastOp <- http.StatusAccepted
			return true
		}
		ReceiveLastOp <- http.StatusTooManyRequests
		return true
	}
	ReceiveLastOp <- http.StatusServiceUnavailable
	return false
}
