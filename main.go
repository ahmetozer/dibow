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





// as util


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
	a := []string{"sh", "df", "gzip", "lsblk"}

	for i, s := range a {
		path, err := exec.LookPath(s)
		if err == nil {
			fmt.Printf("Required program %v %v found at %v\n", i+1, s, path	)
			requiredapps[i] = true
		} else {
			fmt.Printf("Required program %v %v cannot found.\n", i+1, s)
			requiredapps[i] = false
			if i < 2 {
				fmt.Printf("Please install %v and run this program again\n", s)
				os.Exit(3)
			}
		}

	}
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

 func newWebserver(logger *log.Logger) *http.Server {
   router := http.NewServeMux()
   router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
     w.WriteHeader(http.StatusOK)
		 if requiredapps[2] {
			 gzip = "true"
		 } else {
			 gzip = "false"
		 }
		 fmt.Fprintf(w,`<!DOCTYPE html>
		 <html>
		 <head>
		 	<meta name="viewport" content="width=device-width, initial-scale=1.0">

		 <title>HTTP Disk Image Backup</title>
		 <style>
		 	th, td, p, body {
		 		font-family: Arial, Helvetica, sans-serif;
		 		font:16px;
		 	}
		 	h1, h2 {
		 		color: #eff4ff;
		 	}
		 	a {
		 		color: black;
		 	}
		 	table, th, td
		 	{
		 			border: solid 1px #DDD;
		 			border-collapse: collapse;
		 			padding: 2px 3px;
		 			text-align: center;
		 	}
		 	th {
		 			font-weight:bold;
		 	}
		 	body {
		 background-color: #2c3e50;
		 font-size: 16px;
		 font-weight: 400;
		 text-rendering: optimizeLegibility;
		 }

		 div.header {
			 text-align:center;
		 }




		 /*** Table Styles **/

		 .table-fill {
		 background: white;
		 border-radius:3px;
		 border-collapse: collapse;
		 height: 320px;
		 margin: auto;
		 max-width: 600px;
		 padding:5px;
		 width: 100%;
		 box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
		 animation: float 5s infinite;
		 }

		 th {
		 color:#D5DDE5;;
		 background:#1b1e24;
		 border-bottom:4px solid #9ea7af;
		 border-right: 1px solid #343a45;
		 font-size:23px;
		 font-weight: 100;
		 padding:24px;
		 text-align:left;
		 text-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
		 vertical-align:middle;
		 }

		 th:first-child {
		 border-top-left-radius:3px;
		 }

		 th:last-child {
		 border-top-right-radius:3px;
		 border-right:none;
		 }

		 tr {
		 border-top: 1px solid #C1C3D1;
		 border-bottom-: 1px solid #C1C3D1;
		 color:#666B85;
		 font-size:16px;
		 font-weight:normal;
		 text-shadow: 0 1px 1px rgba(256, 256, 256, 0.1);
		 }

		 tr:hover td {
		 background:#4E5066;
		 color:#FFFFFF;
		 border-top: 1px solid #22262e;
		 }

		 tr:first-child {
		 border-top:none;
		 }

		 tr:last-child {
		 border-bottom:none;
		 }

		 tr:nth-child(odd) td {
		 background:#EBEBEB;
		 }

		 tr:nth-child(odd):hover td {
		 background:#4E5066;
		 }

		 tr:last-child td:first-child {
		 border-bottom-left-radius:3px;
		 }

		 tr:last-child td:last-child {
		 border-bottom-right-radius:3px;
		 }

		 td {
		 background:#FFFFFF;
		 padding:20px;
		 text-align:left;
		 vertical-align:middle;
		 font-weight:300;
		 font-size:18px;
		 text-shadow: -1px -1px 1px rgba(0, 0, 0, 0.1);
		 border-right: 1px solid #C1C3D1;
		 }

		 td:last-child {
		 border-right: 0px;
		 }

		 th.text-left {
		 text-align: left;
		 }

		 th.text-center {
		 text-align: center;
		 }

		 th.text-right {
		 text-align: right;
		 }

		 td.text-left {
		 text-align: left;
		 }

		 td.text-center {
		 text-align: center;
		 }

		 td.text-right {
		 text-align: right;
		 }
		 </style>
		 </head>
		 <body>
		 <div class="header">
		 	 <h1>HTTP Disk Image Backup</h1>
			 <h2>Server `+hostname+`</h2>
		 </div>
		 <div id="showData"></div>
		 </body>

		 <script>
		  var gzip = `+gzip+`;
		 	var blockdevices= [
		 			{"name": "fd0", "maj:min": "2:0", "rm": "1", "size": "4K", "ro": "0", "type": "disk", "mountpoint": null},
		 			{"name": "sda", "maj:min": "8:0", "rm": "0", "size": "127G", "ro": "0", "type": "disk", "mountpoint": null,
		 				 "children": [
		 						{"name": "sda1", "maj:min": "8:1", "rm": "0", "size": "127G", "ro": "0", "type": "part", "mountpoint": "/"}
		 				 ]
		 			}
		 	 ]

		 	// EXTRACT VALUE FOR HTML HEADER.
		 	// ('Book ID', 'Book Name', 'Category' and 'Price')
		 	var col = [];
		 	for (var i = 0; i < blockdevices.length; i++) {
		 			for (var key in blockdevices[i]) {
		 					if (col.indexOf(key) === -1) {
		 							col.push(key);
		 					}
		 			}
		 	}

		 	// CREATE DYNAMIC TABLE.
		 	var table = document.createElement("table");
		 	table.classList.add('table-fill');

		 	// CREATE HTML TABLE HEADER ROW USING THE EXTRACTED HEADERS ABOVE.

		 	var tr = table.insertRow(-1);                   // TABLE ROW.

		 	for (var i = 0; i < 7; i++) {
		 			var th = document.createElement("th");      // TABLE HEADER.
		 			th.classList.add('text-left');
		 			th.innerHTML = col[i];
		 			tr.appendChild(th);
		 	}

		 	// ADD JSON DATA TO THE TABLE AS ROWS.
		 	function rows(blockdevices) {
		 	for (var i = 0; i < blockdevices.length; i++) {

		 			tr = table.insertRow(-1);
		 			for (var j = 0; j < 1; j++) {
		 					var tabCell = tr.insertCell(-1);
		 					if (gzip == 1) {
		 						tabCell.innerHTML = "<a href='/image/dev/"+blockdevices[i][col[j]]+"'>"+blockdevices[i][col[j]]+".img</a></br></br><a href='/image.gz/dev/"+blockdevices[i][col[j]]+"'>"+blockdevices[i][col[j]]+".img.gz</a>";
		 					} else {
		 						tabCell.innerHTML = "<a href='/image/dev/"+blockdevices[i][col[j]]+"'>"+blockdevices[i][col[j]]+".img</a>";
		 					}
		 			}
		 			for (var j = 1; j < 7; j++) {
		 					var tabCell = tr.insertCell(-1);
		 					tabCell.innerHTML = blockdevices[i][col[j]];
		 			}

		 			for (var j = 7; j <= 7; j++) {
		 				if ( blockdevices[i][col[j]] != undefined ) {
		 					rows(blockdevices[i][col[j]])
		 				}
		 			}

		 	}
		 }
		  rows(blockdevices)
		 	// FINALLY ADD THE NEWLY CREATED TABLE WITH JSON DATA TO A CONTAINER.
		 	var divContainer = document.getElementById("showData");
		 	divContainer.innerHTML = "";
		 	divContainer.appendChild(table);

		 </script>
		 </html>`)
   })
	 router.HandleFunc("/disk/", func(w http.ResponseWriter, r *http.Request) {
     w.WriteHeader(http.StatusOK)
		 fmt.Fprintf(w, "Selected disk %s", r.URL.Path[5:])
   })

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

	 router.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
			 t := time.Now()
		 filename := "backup-"+hostname+"-"+t.Format(time.RFC3339)+".img"
		 w.Header().Set("Content-Disposition", "attachment; filename="+filename)


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
