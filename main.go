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
	beats, port                  int
	item                         string
	localCircuit1, localCircuit2 Circuit
	heartbeat                    *time.Ticker
	router                       *gin.Engine
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

	// Make sure we are still alive.
	router.GET("/stats/:pgm", GetCircuitStats)

	// These are the services we will be listening for.
	router.POST("/add/:word", ReceiveWrapper)
	// Get the number of heartbeats put out by the application (also in real-time).
	router.GET("/beats", GetHeartbeatCount)
	// Make sure we are still alive.
	router.GET("/ping", PingTheAPI)
} // func

// ReceiveWrapper ...
func ReceiveWrapper(c *gin.Context) {
	if localCircuit1.Switch.LastOperation != http.StatusServiceUnavailable {
		localCircuit1.Conductor.Fill(c.Param("word"))
	}
	content := gin.H{"payload": len(localCircuit1.Conductor.Channel)}
	c.JSON(localCircuit1.Switch.LastOperation, content)
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

// GetCircuitStats ...
func GetCircuitStats(c *gin.Context) {
	var content map[string]interface{}
	type status struct {
		Name          string `json:"name"`
		Depth         int    `json:"depth"`
		Switch        int    `json:"switch"`
		LastOperation int    `json:"last_operation"`
	}
	switch c.Param("pgm") {
	case "1":
		content = gin.H{"payload": status{Name: localCircuit1.Conductor.Name, Depth: len(localCircuit1.Conductor.Channel), LastOperation: localCircuit1.Conductor.LastOperation}}
	case "1h":
		content = gin.H{"payload": status{Name: localCircuit1.Switch.Name, Depth: len(localCircuit1.Switch.Channel), LastOperation: localCircuit1.Switch.LastOperation}}
	case "2":
		content = gin.H{"payload": status{Name: localCircuit2.Conductor.Name, Depth: len(localCircuit2.Conductor.Channel), Switch: len(localCircuit2.Switch.Channel), LastOperation: localCircuit2.Conductor.LastOperation}}
	case "2h":
		content = gin.H{"payload": status{Name: localCircuit2.Switch.Name, Depth: len(localCircuit2.Switch.Channel), LastOperation: localCircuit2.Switch.LastOperation}}
	} // switch
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
				fmt.Printf(`{"date":"%s","heartbeat":"%d","receiver1 depth":"%d","hold1 depth":"%d","status":"%d"}`+"\n", time.Now().UTC(), beats, len(localCircuit1.Conductor.Channel), len(localCircuit1.Switch.Channel), localCircuit1.Conductor.LastOperation)
				fmt.Printf(`{"date":"%s","heartbeat":"%d","receiver2 depth":"%d","hold2 depth":"%d","status":"%d"}`+"\n", time.Now().UTC(), beats, len(localCircuit2.Conductor.Channel), len(localCircuit2.Switch.Channel), localCircuit2.Conductor.LastOperation)
			} // select
		} // for
	}() // go func

	// Spoken Word App
	go func() {
		fmt.Printf("Starting Spoken Word App (1)...\n")
		localCircuit1 = Circuit{}
		localCircuit1 = *localCircuit1.New()
		//localCircuit1.Conductor.LastOperation = http.StatusAccepted
		for {
			// Watch for stuff to happen.
			select {
			case <-localCircuit1.Conductor.Check():
				fmt.Printf("Speaking Word: %d\n", len(localCircuit1.Conductor.Channel))
				if localCircuit2.Conductor.LastOperation != http.StatusServiceUnavailable {
					localCircuit1.Send.Fill(localCircuit1.Conductor.Deplete())
					// Slow things down if there are too many requests.
				} else {
					if localCircuit1.Switch.LastOperation != http.StatusServiceUnavailable {
						localCircuit1.Switch.Fill(localCircuit1.Conductor.Deplete())
					}
					localCircuit1.Conductor.LastOperation = localCircuit1.Switch.LastOperation
				}
			case <-localCircuit1.Switch.Check():
				fmt.Printf("Speaking Word [HOLD]: %d\n", len(localCircuit1.Switch.Channel))
			} // select
		} // for
	}() // go func

	// Written Word App
	go func() {
		fmt.Printf("Starting Written Word App (2)...\n")
		localCircuit2 = Circuit{}
		localCircuit2 = *localCircuit2.New()
		localCircuit1.Send = &localCircuit2.Conductor
		//localCircuit2.Conductor.LastOperation = http.StatusAccepted
		for {
			// Watch for stuff to happen.
			select {
			case <-localCircuit2.Conductor.Check():
				fmt.Printf("Writing Word: %d\n", len(localCircuit2.Conductor.Channel))
			case <-localCircuit2.Switch.Check():
				fmt.Printf("Written Word [HOLD]: %d\n", len(localCircuit2.Switch.Channel))
			} // select
		} // for
	}() // go func

	// After the 'go func' is dispatched, start the server and listen on the
	// specified port.
	fmt.Printf("%s gosp3 message=engine started, ready on port %d\n", time.Now().UTC(), port)
	router.Run(":7718")
} // func
