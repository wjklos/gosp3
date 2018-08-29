package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GetHashNotifier ...
func GetHashNotifier() chan bool {
	return hashNotifier
}

// NotifyHashList ...
func NotifyHashList() {
	// Since this is only a "call to action" channel, it only needs one call.
	// If there is already a message in it, then someone else made that call.
	if len(hashNotifier) < cap(hashNotifier) {
		hashNotifier <- true
	}
}

// FillHashList ...
func FillHashList(item string) bool {
	// Make sure there is room in the list before adding any thing to it.
	if len(hash) < cap(hash) {
		hash <- item
		NotifyHashList()
		return true
	}
	return false
}

// DepleteHashList ...
func DepleteHashList() bool {
	var item string
	// Spin through the channel looking to process anything in it.
	z := len(hash)
	for i := 0; i < z; i++ {
		item = <-hash
		h := md5.New()
		h.Write([]byte(item))
		fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil)))
	} // for
	return true
} // func
