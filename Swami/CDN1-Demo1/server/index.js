const express = require('express')
const http = require('http')
const axios = require('axios');
const cors = require('cors');
const _ = require('lodash');

var app = express()
app.use(cors());
var port = process.env.PORT || 8080

// When user posts request for SLA info,
// send get request to Swami
// return information received from swami
app.get('/getSLA', (req, res) => {
  var sessionID = uuid();

  axios.get(`http://18.216.195.86:8081/getpaths?source=${req.query.source}&dest=${req.query.destination}&sid=${sessionID}`).then((response) => {
    res.send(response.data)
  }).catch(error => {
    console.log(error)
  })

})

app.post('/setSLA', (req, res) => {
  sessionID = req.query.sid
  pathSelected = req.query.path

  axios.post(`http://18.216.195.86:8081/setpath?sid=${sessionID}&path=${pathSelected}`).then((response) => {
    res.send('success')
  }).catch((error) => {
    // on error
  })
})

app.post('/resettest', (req, res) => {
  axios.post(`http://18.216.195.86:8081/resettest`).then((response) => {
    res.send('success')
  }).catch((error) => {
    // on error
  })
})

// app.use('/', express.static('public'))
let httpServer = http.Server(app);
httpServer.listen(port);
app.use(express.static('public'));

// app.listen(port, () => {
//   console.log(`Started on port ${port}`)
// })

function uuid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
    }

    return s4() + s4() + s4() + s4() + s4() + s4() + s4() + s4();
}
