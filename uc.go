package main

// import (
// 	"fmt"
// 	"strings"
// )

// // GetUCNotifier ...
// func GetUCNotifier() chan bool {
// 	return ucNotifier
// }

// // NotifyUCList ...
// func NotifyUCList() {
// 	// Since this is only a "call to action" channel, it only needs one call.
// 	// If there is already a message in it, then someone else made that call.
// 	if len(ucNotifier) < cap(ucNotifier) {
// 		ucNotifier <- true
// 	}
// }

// // FillUCList ...
// func FillUCList(item string) bool {
// 	// Make sure there is room in the list before adding any thing to it.
// 	if len(uc) < cap(uc) {
// 		uc <- item
// 		NotifyUCList()
// 		return true
// 	}
// 	return false
// }

// // DepleteUCList ...
// func DepleteUCList() bool {
// 	var item string
// 	// Spin through the channel looking to process anything in it.
// 	z := len(uc)
// 	for i := 0; i < z; i++ {
// 		item = <-uc
// 		fmt.Printf("%s\n", strings.ToUpper(item))
// 	} // for
// 	return true
// } // func
