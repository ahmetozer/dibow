package main

import (
	"io"
	"net/http"
	"os/exec"
	"log"
	"fmt"
	"os"
	"flag"
	"context"
	"os/signal"
	"time"
)

var (
	BUF_LEN = 1024
	listenAddr string
	requiredapps[4]bool
	hostname string
	gzip string
)





// pass CMD output to HTTP


func writeCmdOutput(res http.ResponseWriter, pipeReader *io.PipeReader) {
	buffer := make([]byte, BUF_LEN)
	for {
		n, err := pipeReader.Read(buffer)
		if err != nil {
			pipeReader.Close()
			break
		}

		data := buffer[0:n]
		res.Write(data)
		if f, ok := res.(http.Flusher); ok {
			f.Flush()
		}
		//reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}



func main() {

	// Requred programs for the backup system on Linux enviroment

	required_programs := []string{"sh", "df", "gzip", "lsblk"}

	for i, s := range required_programs {
		// Get path of required programs
		path, err := exec.LookPath(s)
		if err == nil {
			fmt.Printf("Required program %v %v found at %v\n", i+1, s, path	)
			requiredapps[i] = true	//save to array
		} else {
			fmt.Printf("Required program %v %v cannot found.\n", i+1, s)
			requiredapps[i] = false
			if i < 2 {							//sh and df is must required. If is not found in software than exit.
				fmt.Printf("Please install %v and run this program again\n", s)
				os.Exit(3)
			}
		}

	}

	//	Get Hostname of Server
	var err error
	hostname, err = os.Hostname()
		if err != nil {
			panic(err)
		}

	fmt.Println("hostname:", hostname)

	flag.StringVar(&listenAddr, "listen-addr", ":80", "server listen address")
  flag.Parse()

  logger := log.New(os.Stdout, "http: ", log.LstdFlags)

  done := make(chan bool, 1)
  quit := make(chan os.Signal, 1)

  signal.Notify(quit, os.Interrupt)

  server := newWebserver(logger)
  go gracefullShutdown(server, logger, quit, done)

  logger.Println("Server is ready to handle requests at", listenAddr)
  if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
  }

  <-done
  logger.Println("Server stopped")
 }

	//Grace Full Shutdown
 func gracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
   <-quit
   logger.Println("Server is shutting down...")

   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

   server.SetKeepAlivesEnabled(false)
   if err := server.Shutdown(ctx); err != nil {
     logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
   }
   close(done)
 }


	
