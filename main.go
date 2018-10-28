package main

import (
	"fmt"
	"net/http"
	"time"

	// We can do this natively just as easily, but this framework makes
	// the examples a bit more clear.
	"github.com/gin-gonic/gin"
)

const (
	maxMessages = 10
)

var (
	beats, port int
	item        string

	heartbeat  *time.Ticker
	router     *gin.Engine
)

// Do all of this stuff first.
func init() {

	// In this example, we will hard code the port.  Later the environment
	// will dictate.
	port = 7718
	// Set up the heartbeat ticker.
	heartbeat = time.NewTicker(60 * time.Second)

	// Setup the service router.
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// These are the services we will be listening for.
	router.POST("/add/:word", ReceiveWrapper)
	// Get the number of heartbeats put out by the application (also in real-time).
	router.GET("/beats", GetHeartbeatCount)
	// Make sure we are still alive.
	router.GET("/ping", PingTheAPI)

} // func

// ReceiveWrapper ...
func ReceiveWrapper(c *gin.Context) {
	FillReceiveList(c.Param("word"))
	content := gin.H{"payload": len(ReceiveList)}
	c.JSON(<-ReceiveLastOp, content)
}

// GetHeartbeatCount sends the number of times the heartbeat ticker has
// fired since the program started.
func GetHeartbeatCount(c *gin.Context) {
	content := gin.H{"payload": beats}
	c.JSON(http.StatusOK, content)
}

// PingTheAPI lets the caller know we are alive.
func PingTheAPI(c *gin.Context) {
	content := gin.H{"payload": "pong"}
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
			case <-GetReceiveNotifier():
				fmt.Printf("Adding Word: %d\n", len(ReceiveList))
			// When the Heartbeat ticker is fired, execute this.
			case <-heartbeat.C:
				beats++
				fmt.Printf(`{"date":"%s","app":"gosp2","msgtype":"info","heartbeat":"%d"}`+"\n", time.Now().UTC(), beats)
			} // select
		} // for
	}() // go func

	// After the 'go func' is dispatched, start the server and listen on the
	// specified port.
	fmt.Printf("%s gosp2 msgtype=info message=engine started, ready on port %d\n", time.Now().UTC(), port)
	router.Run(":7718")
} // func
