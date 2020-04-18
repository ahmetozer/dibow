package diserver

// index.html

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // return 200
	if requiredapps[2] {         // Set gzip variable for javascript
		gzip = "true"
	} else {
		gzip = "false"
	}
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html>
	<head>
	 <meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>Dibow</title>
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
		<h1>Dibow</h1>
		<h2>Server `+hostname+`</h2>
		<button onclick="load_disk()">Reload Disks</button>
	</div>
	<div id="showData"></div>
	</body>

	<script>
	 var gzip = `+gzip+`;
	 var getJSON = function(url, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url, true);
    xhr.responseType = 'json';
    xhr.onload = function() {
      var status = xhr.status;
      if (status === 200) {
        callback(null, xhr.response);
      } else {
        callback(status, xhr.response);
      }
    };
    xhr.send();
};

	function load_disk() {
	let url = 'lsblk.json';
	fetch(url)
	.then(res => res.json())
	.then((out) => {
		//blockdevices = out;
		 var blockdevices= out.blockdevices
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
		rows(blockdevices);

		 // FINALLY ADD THE NEWLY CREATED TABLE WITH JSON DATA TO A CONTAINER.
		 var divContainer = document.getElementById("showData");
		 divContainer.innerHTML = "";
		 divContainer.appendChild(table);
	})
	.catch(err => { throw err });
}
load_disk()

	</script>
	</html>`)
}
