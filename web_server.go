package main

import (
	"io"
	"net/http"
	"os/exec"
	"log"
	"fmt"
	"time"
)

func newWebserver(logger *log.Logger) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	 index(w, r)
	})


	// Test function
	router.HandleFunc("/disk/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Selected disk %s", r.URL.Path[5:])
	})

	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/lsblk.json", func(w http.ResponseWriter, r *http.Request) {
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
	 }	else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "lsblk not installed")
	 }
	})

	// Get SYSTEM IMAGE
	router.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
		filename := "backup-"+hostname+"-"+t.Format(time.RFC3339)+".img"
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		// Set HTTP header befare Transfering data.
	 w.Header().Set("Transfer-Encoding", "chunked")
	 cmd := exec.Command("dd",fmt.Sprintf("if=%s", r.URL.Path[6:]))
	 //cmd := exec.Command("bash", "run.sh")
	 pipeReader, pipeWriter := io.Pipe()
	 cmd.Stdout = pipeWriter
	 cmd.Stderr = pipeWriter
	 go writeCmdOutput(w, pipeReader)
	 cmd.Run()
	 pipeWriter.Close()
	})

	// Get SYSTEM IMAGE with gzipped to save bandwith
	router.HandleFunc("/image.gz/", func(w http.ResponseWriter, r *http.Request) {
	 if requiredapps[2] == true {
		 // Get time and set headers
		 t := time.Now()
		 filename := "backup-"+hostname+"-"+t.Format(time.RFC3339)+".img.gz"
		 w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		 w.Header().Set("Transfer-Encoding", "chunked")

		 // run Command
		 cmd := exec.Command("dd",fmt.Sprintf("if=%s", r.URL.Path[9:]))
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


	})

	return &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 10 * time.Second,
		//IdleTimeout:  15 * time.Second,
	}
}
