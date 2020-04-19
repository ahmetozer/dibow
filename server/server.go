package diserver

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

var (
	BUF_LEN      = 1024
	listenAddr   string
	requiredapps [4]bool
	hostname     string
	gzip         string
)

var (
	host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	validFrom  = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	rsaBits    = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	ed25519Key = flag.Bool("ed25519", false, "Generate an Ed25519 key")
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}



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


//	Server
func Server(args []string) {


	// Checking required applications.
	required_programs := []string{"sh", "df", "gzip", "lsblk"}

	for i, s := range required_programs {
		// Get path of required programs
		path, err := exec.LookPath(s)
		if err == nil {
			fmt.Printf("Required program %v %v found at %v\n", i+1, s, path)
			requiredapps[i] = true //save to array
			} else {
				fmt.Printf("Required program %v %v cannot found.\n", i+1, s)
				requiredapps[i] = false
				if i < 2 { //sh and df is must required. If is not found in software than exit.
					fmt.Printf("Please install %v and run this program again\n", s)
					os.Exit(3)
				}
			}
		}

		flag := flag.NewFlagSet("Server", flag.ContinueOnError)

		//	Get Hostname of the Server
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			panic(err)
		}

		//	Print hostname
		fmt.Println("hostname:", hostname)
		flag.StringVar(&listenAddr, "listen-addr", ":443", "server listen address")
		flag.Parse(args)
		//	Parse arguments
		//	Log to file
		logger := log.New(os.Stdout, "https: ", log.LstdFlags)
		done := make(chan bool, 1)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		//	Generate SSL Certificate For the webserver
		ssl_cert_generate()

		//	Start the web server. // web_server.go
		server := newWebserver(logger)
		go gracefullShutdown(server, logger, quit, done)
		logger.Println("Server is ready to handle requests at", listenAddr)
		if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
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
