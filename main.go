package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	// We can do this natively just as easily, but this framework makes
	// the examples a bit more clear.
	"github.com/gin-gonic/gin"
)

const (
	maxMessages = 100
)

var (
	beats, port              int
	lc, uc, hash             chan string
	item                     string
	ucNotifier, hashNotifier chan bool
	heartbeat                *time.Ticker
	router                   *gin.Engine
)

// Do all of this stuff first.
func init() {
	ucNotifier = make(chan bool, 1)

	lc = make(chan string, maxMessages)
	uc = make(chan string, maxMessages)
	hash = make(chan string, maxMessages)
	// In this example, we will hard code the port.  Later the environment
	// will dictate.
	port = 7718
	// Set up the heartbeat ticker.
	heartbeat = time.NewTicker(60 * time.Second)

	// Setup the service router.
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// These are the services we will be listening for.
	router.POST("/add/:word", AddToLowerCaseQ)
	router.GET("/beats", GetHeartbeatCount)
	router.GET("/ping", PingTheAPI)

} // func

// AddToLowerCaseQ ...
func AddToLowerCaseQ(c *gin.Context) {
	lc <- c.Param("word")
	content := gin.H{"payload": "Accepted: " + fmt.Sprintf("l: %d, U: %d, #: %d", len(lc), len(uc), len(hash))}
	c.JSON(http.StatusOK, content)
}

// GetHeartbeatCount sends the number of times the heartbeat ticker has
// fired since the program started.
func GetHeartbeatCount(c *gin.Context) {
	content := gin.H{"payload": beats}
	c.JSON(http.StatusOK, content)
}

// PingTheAPI lets the caller know we are alive.
func PingTheAPI(c *gin.Context) {
	content := gin.H{"payload": "PONG"}
	c.JSON(http.StatusOK, content)
}

// Manage the processes.
func main() {
	// Dispatch a process into the background.
	go func() {
		// Now run it forever.
		for {
			// Watch for stuff to happen.
			select {
			// When the Heartbeat ticker is fired, execute this.
			case <-heartbeat.C:
				beats++
				fmt.Printf("bump,Bump... @ %s\n", time.Now().UTC())
			} // select
		} // for
	}() // go func

	// Make lowercase.
	go func() {
		for {
			if len(lc) > 0 && (len(lc) < cap(lc)) {
				item = strings.ToLower(<-lc)
				fmt.Printf("%s\n", item)
				uc <- item
			} // if
		} // for
	}() // go func

	// Make UPPERCASE.
	go func() {
		for {
			if len(uc) > 0 && (len(uc) < cap(uc)) {
				item = strings.ToUpper(<-uc)
				fmt.Printf("%s\n", item)
				hash <- item
			} // if
		} // for
	}() // go func

	// Make hash.
	go func() {
		for {
			if len(hash) > 0 && (len(hash) < cap(hash)) {
				h := md5.New()
				h.Write([]byte(<-hash))
				fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil)))
			} // if
		}
	}()

	// After the 'go func' is dispatched, start the server and listen on the
	// specified port.
	fmt.Printf("ready on port %d\n", port)
	router.Run(":" + strconv.Itoa(port))
} // func
