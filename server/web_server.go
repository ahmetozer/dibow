package diserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"os"
	"time"
	"math/rand"
	"encoding/base64"
	"strings"
	//"bytes"
)

var (
	password = random_password()
)
func random_password() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	digits + specials
	length := 16
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	for i := len(buf) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}
	str := string(buf)

	return str
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != "root" || pair[1] != password {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func HttpClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}
// pass CMD output to HTTP


// Web server to handle disk operation requests
func newWebserver(logger *log.Logger) *http.Server {

	//Print Root user Password
	fmt.Printf("Random password: %s\n", password)

	// Crearte New HTTP Router
	router := http.NewServeMux()

	// Index Handler
	router.HandleFunc("/", index)

	// Disk Test function
	router.HandleFunc("/disk/", use(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Selected disk %s", r.URL.Path[5:])
	}, basicAuth) )

	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/lsblk.json", use(func(w http.ResponseWriter, r *http.Request) {

			//	Check lsblk application
			if requiredapps[3] == true {
				// return success
				w.WriteHeader(http.StatusOK)
				cmd := exec.Command("lsblk", "-J")
				// Organize pipelines
				pipeIn, pipeWriter := io.Pipe()
				cmd.Stdout = pipeWriter
				cmd.Stderr = pipeWriter
				// Pass to web output
				go writeCmdOutput(w, pipeIn)

				// Run commands
				cmd.Run()
				pipeWriter.Close()

				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "lsblk not installed")
				}
	}, basicAuth))

	// Get SYSTEM IMAGE
	router.HandleFunc("/image/", use(func(w http.ResponseWriter, r *http.Request) {
	// Check device path.
	// This also prevent if any vulnerable is avaible at url
	if fileExists(r.URL.Path[6:]) {
		t := time.Now()
		filename := "backup-" + hostname + "-" + t.Format(time.RFC3339) + ".img"
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		// Set HTTP header befare Transfering data.
		w.Header().Set("Transfer-Encoding", "chunked")

		ddCommand := exec.Command("dd", fmt.Sprintf("if=%s", r.URL.Path[6:]))
		pipeIn, pipeWriter := io.Pipe()
		ddCommand.Stdout = pipeWriter
		ddCommand.Stderr = pipeWriter

		go writeCmdOutput(w, pipeIn)

		ddCommand.Run()
		pipeWriter.Close()
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("Disk %s is does not exist. This query requests by %s\n", r.URL.Path[6:], HttpClientIP(r) )
	}
	}, basicAuth))

	// Get SYSTEM IMAGE with gzipped to save bandwith
	router.HandleFunc("/image.gz/", use(func(w http.ResponseWriter, r *http.Request) {

		if requiredapps[2] == true { // Check gzip on system

			if fileExists(r.URL.Path[9:]) { // Check device path. // This also prevent if any vulnerable is avaible at
				t := time.Now()
				filename := "backup-" + hostname + "-" + t.Format(time.RFC3339) + ".img.gz"
				w.Header().Set("Content-Disposition", "attachment; filename="+filename)
				// Set HTTP header befare Transfering data.
				w.Header().Set("Transfer-Encoding", "chunked")

				ddCommand := exec.Command("sh", "-c", "dd "+fmt.Sprintf("if=%s", r.URL.Path[9:])+" | gzip -")
				pipeIn, pipeWriter := io.Pipe()
				ddCommand.Stdout = pipeWriter
				ddCommand.Stderr = pipeWriter

				go writeCmdOutput(w, pipeIn)

				ddCommand.Run()
				pipeWriter.Close()



			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Printf("Disk %s is does not exist. This query requests by %s\n", r.URL.Path[9:], HttpClientIP(r) )
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "gzip cannot found\n")
		}

		}, basicAuth))

		return &http.Server{
			Addr:     listenAddr,
			Handler:  router,
			ErrorLog: logger,
			//ReadTimeout:  5 * time.Second,
			//WriteTimeout: 10 * time.Second,
			//IdleTimeout:  15 * time.Second,
		}
}
