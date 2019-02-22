const express = require('express')
const http = require('http')
const axios = require('axios');
const cors = require('cors');
const _ = require('lodash');
const fs = require('fs');

var app = express()
app.use(cors());
var port = process.env.PORT || 8080
// const ORACLE_ADDR = "localhost"
const ORACLE_ADDR = "18.216.195.86"
// const ORACLE_ADDR = "155.246.45.43"
const SESSION_ID_TEST = "1001100110011001"

var lastUsedProcessNum = 0
var currentProcesses = {
  10000: {
    myPath: "/index.html",
    res: "res",
    onCompletion: function() {
    console.log("onCompletion")
    }
  },
}

// When user posts request for SLA info,
// send get request to Swami
// return information received from swami
app.get('/getSLA', (req, res) => {
  var sessionID = uuid();

  axios.get(`http://${ORACLE_ADDR}:8081/getpaths?source=${req.query.source}&dest=${req.query.destination}&sid=${sessionID}`).then((response) => {
    res.send(response.data)
  }).catch(error => {
    console.log(error)
  })

})

app.post('/setSLA', (req, res) => {
  sessionID = req.query.sid
  pathSelected = req.query.path

  axios.post(`http://${ORACLE_ADDR}:8081/setpath?sid=${sessionID}&path=${pathSelected}`).then((response) => {
    res.send(response.data)
  }).catch((error) => {
    // on error
  })
})

app.post('/resettest', (req, res) => {
  axios.post(`http://${ORACLE_ADDR}:8081/resettest`).then((response) => {
    res.send(response.data)
  }).catch((error) => {
    // on error
  })
})

app.get('/test.mp4', (req, res) => {
  // var myPath = "public"+req.path
  // if (fs.existsSync(myPath)) {
  //   // var content;
  //   // fs.readFile(myPath, function(err, data) {
  //   //   if (err) {
  //   //     throw err;
  //   //   }
  //   //
  //   //   content = data;
  //   //   res.send(content);
  //   // })
  //   obtainFile(myPath, res)
  //   console.log(myPath+" does exist")
  // } else {
  //   obtainFile(myPath, res)
  //   console.log(myPath+" does not exist")
  // }
  obtainFile("public/test.mp4", res)
})

function obtainFile(myPath, res) {
  let currentProcessNum = lastUsedProcessNum+1;
  for (p in currentProcesses) {
    if (p.myPath == myPath) {
      fmt.Println("Currently retreiving file. This request will be ignored. No response will be given")
      return
    }
  }

  currentProcesses[currentProcessNum] = {
    processNum: currentProcessNum,
    myPath: myPath,
    res: res,
    start: function() {
      console.log("STARTING")
      axios.post(`http://${ORACLE_ADDR}:8081/getfile?file=${this.myPath}&sid=${SESSION_ID_TEST}`).then((response) => {
        console.log("RESPONSE:", response.data)
        this.onCompletion(response)
      }).catch((error) => {
        console.log(error)
      })
    },
    onCompletion: function() {
      console.log("COMPLETION")
      var content;
      var response = this.res
      fs.readFile(this.myPath, function(err, data, res=response) {
        if (err) {
          throw err;
        }

        content = data;
        res.send(content);
        delete currentProcesses[this.processNum]
      })

      console.log(myPath+" does exist")
      delete currentProcesses[this.currentProcessNum]
    }
  }

  currentProcesses[currentProcessNum].start()
}

app.use('/', express.static('public'))
let httpServer = http.Server(app);
httpServer.listen(port);
// app.use(express.static('public'));

// app.listen(port, () => {
//   console.log(`Started on port ${port}`)
// })

function uuid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
    }

    return s4() + s4() + s4() + s4() + s4() + s4() + s4() + s4();
}


function WebSocket() {
            
  if ("WebSocket" in window) {
     alert("WebSocket is supported by your Browser!");
     
     // Let us open a web socket
     var ws = new WebSocket("ws://localhost:8000/echo");

     ws.onopen = function() {
        
        // Web Socket is connected, send data using send()
        ws.send("Want to send a message?");
        alert("Sent message");
     };

     ws.onmessage = function (evt) { 
        var received_msg = evt.data;
        alert("You recieved a message");
     };

     ws.onclose = function() { 
        
        // websocket is closed.
        alert("You closed the connection"); 
     };
  } else {
    
     // The browser doesn't support WebSocket
     alert("WebSocket NOT supported by your Browser!");
  }
}