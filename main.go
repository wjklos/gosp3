package main

import (
	"fmt"
	"net/http"
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
	beats, port                                     int
	item                                            string
	origAddress, lcAddress, titleAddress, ucAddress []string
	lc, title, uc                                   chan string
	heartbeat                                       *time.Ticker
	router                                          *gin.Engine
)

// Do all of this stuff first.
func init() {
	lc = make(chan string, maxMessages)
	uc = make(chan string, maxMessages)
	title = make(chan string, maxMessages)
	// In this example, we will hard code the port.  Later the environment
	// will dictate.
	port = 7718
	// Set up the heartbeat ticker.
	heartbeat = time.NewTicker(60 * time.Second)

	// Setup the service router.
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// These are the services we will be listening for.
	// Insert words into the process. See gettysburg.sh for example.
	router.POST("/add/:word", AddToLowerCaseQ)
	// Show the total words processed by each stage (in real-time).
	router.GET("/totals", GetTotals)
	// Get the number of heartbeats put out by the application (also in real-time).
	router.GET("/beats", GetHeartbeatCount)
	// Show the results of the lowercase process.
	router.GET("/lc", GetLCAddress)
	// Show the results of the TitleCase process.
	router.GET("/title", GetTitleAddress)
	// Show the results of the UPPERCASE process.
	router.GET("/uc", GetUCAddress)
	// Make sure we are still alive.
	router.GET("/ping", PingTheAPI)

} // func

// AddToLowerCaseQ ...
func AddToLowerCaseQ(c *gin.Context) {
	lc <- c.Param("word")
	origAddress = append(origAddress, c.Param("word"))
	content := gin.H{"payload": "Accepted: " + fmt.Sprintf("lc: %d, Tc: %d, UC: %d, #: %d", len(lc), len(title), len(uc), len(origAddress))}
	c.JSON(http.StatusOK, content)
}

// GetTotals ...
func GetTotals(c *gin.Context) {
	content := gin.H{"payload": fmt.Sprintf("lc: %d, Tc: %d, UC: %d, #: %d", len(lcAddress), len(titleAddress), len(ucAddress), len(origAddress))}
	c.JSON(http.StatusOK, content)
}

// GetLCAddress ...
func GetLCAddress(c *gin.Context) {
	content := gin.H{"payload": strings.Join(lcAddress, " ")}
	c.JSON(http.StatusOK, content)
}

// GetTitleAddress ...
func GetTitleAddress(c *gin.Context) {
	content := gin.H{"payload": strings.Join(titleAddress, " ")}
	c.JSON(http.StatusOK, content)
}

// GetUCAddress ...
func GetUCAddress(c *gin.Context) {
	content := gin.H{"payload": strings.Join(ucAddress, " ")}
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
			// When the Heartbeat ticker is fired, execute this.
			case <-heartbeat.C:
				beats++
				fmt.Printf(`{"date":"%s","app":"gosp2","msgtype":"info","heartbeat":"%d"}`+"\n", time.Now().UTC(), beats)
				//fmt.Printf("%s gosp2 msgtype=info heartbeat=%d\n", time.Now().UTC(), beats)
				//fmt.Printf("bump,Bump... @ %s\n", time.Now().UTC())
			} // select
		} // for
	}() // go func

	// Make lowercase.
	go func() {
		for {
			select {
			case item := <-lc:
				item = strings.ToLower(item)
				lcAddress = append(lcAddress, item)
				time.Sleep(time.Millisecond * 50)
				title <- item
			}
		} // for
	}() // go func

	// // Make Title.
	go func() {
		for {
			select {
			case item = <-title:
				item = strings.Title(item)
				titleAddress = append(titleAddress, item)
				time.Sleep(time.Millisecond * 100)
				uc <- item
			}
		}
	}()

	// // Make UPPERCASE.
	go func() {
		for {
			select {
			case item = <-uc:
				item = strings.ToUpper(item)
				ucAddress = append(ucAddress, item)
			}
		} // for
	}() // go func

	// After the 'go func' is dispatched, start the server and listen on the
	// specified port.
	fmt.Printf("%s gosp2 msgtype=info message=engine started, ready on port %d\n", time.Now().UTC(), port)
	router.Run(":7718")
	//router.Run(":" + strconv.Itoa(port))
} // func
