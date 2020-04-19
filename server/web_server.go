package diserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"
	"math/rand"
	"encoding/base64"
	"strings"

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
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}



func newWebserver(logger *log.Logger) *http.Server {

	fmt.Printf("Random password: %s\n", password)
	router := http.NewServeMux()

	router.HandleFunc("/", index)

	// Test function
	router.HandleFunc("/disk/", use(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Selected disk %s", r.URL.Path[5:])
	}, basicAuth) )

	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/lsblk.json", use(func(w http.ResponseWriter, r *http.Request) {
		if requiredapps[2] == true {
			w.WriteHeader(http.StatusOK)
			cmd := exec.Command("lsblk", "-J")
			//cmd := exec.Command("bash", "run.sh")
			pipeReader, pipeWriter := io.Pipe()
			cmd.Stdout = pipeWriter
			cmd.Stderr = pipeWriter
			go writeCmdOutput(w, pipeReader)
			cmd.Run()
			pipeWriter.Close()
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "lsblk not installed")
		}
	}, basicAuth))

	// Get SYSTEM IMAGE
	router.HandleFunc("/image/", use(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		filename := "backup-" + hostname + "-" + t.Format(time.RFC3339) + ".img"
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		// Set HTTP header befare Transfering data.
		w.Header().Set("Transfer-Encoding", "chunked")
		cmd := exec.Command("dd", fmt.Sprintf("if=%s", r.URL.Path[6:]))
		//cmd := exec.Command("bash", "run.sh")
		pipeReader, pipeWriter := io.Pipe()
		cmd.Stdout = pipeWriter
		cmd.Stderr = pipeWriter
		go writeCmdOutput(w, pipeReader)
		cmd.Run()
		pipeWriter.Close()
	}, basicAuth))

	// Get SYSTEM IMAGE with gzipped to save bandwith
	router.HandleFunc("/image.gz/", use(func(w http.ResponseWriter, r *http.Request) {
		if requiredapps[2] == true {
			// Get time and set headers
			t := time.Now()
			filename := "backup-" + hostname + "-" + t.Format(time.RFC3339) + ".img.gz"
			w.Header().Set("Content-Disposition", "attachment; filename="+filename)
			w.Header().Set("Transfer-Encoding", "chunked")

			// run Command
			cmd := exec.Command("dd", fmt.Sprintf("if=%s", r.URL.Path[9:]))
			//c2 := exec.Command("gzip ", "-1", "-")
			//cmd := exec.Command("bash", "run.sh")
			//pipeReader, pipeWriter := io.Pipe()
			pipeReader1, pipeWriter1 := io.Pipe()
			cmd.Stdout = pipeWriter1
			cmd.Stderr = pipeWriter1

			c2 := exec.Command("gzip ", "-1", "-")
			//cmd := exec.Command("bash", "run.sh")

			pipeReader2, pipeWriter2 := io.Pipe()

			c2.Stdin = pipeReader1
			c2.Stdout = pipeWriter2
			c2.Stderr = pipeWriter2

			go writeCmdOutput(w, pipeReader2)
			//cmd.Run()
			cmd.Start()
			c2.Start()
			cmd.Wait()
			c2.Wait()
			pipeWriter1.Close()
			pipeWriter2.Close()

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
