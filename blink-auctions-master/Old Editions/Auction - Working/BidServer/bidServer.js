const HTTPS_PORT = 443;

const nodeStatic = require('node-static');
const https = require('https');
const socketIO = require('socket.io');
const fs = require('fs');
const os = require('os');

/************  SERVER SETUP *************/

const certOptions = {
  key: fs.readFileSync('certs/key.pem'),
  cert: fs.readFileSync('certs/cert.pem')
}

var fileServer = new(nodeStatic.Server)();
var app = https.createServer(certOptions, function(req, res) {
  fileServer.serve(req, res);
}).listen(HTTPS_PORT);

var io = socketIO.listen(app);
console.log("Connected.");

io.sockets.on('connection', function(socket) {

  socket.on('connect to bid', function(userID, roomName) {
    console.log("Connect happened");
    sockets[userID] = socket;
    updateBidders(roomName);
  });

  socket.on('place bid', function(bid) {
    var roomName = bid.roomName;

    placeBid(roomName, bid, socket);
    addBidToUserLog(bid);
  });

});

/******* SETUP MAIN SERVER CONNECTION *********/

var io_client = require('socket.io-client');
var mySocket = io_client.connect("http://bid.blinkcdn.com");
mySocket.emit('connect service', "https://bidserver.blinkcdn.com", "bid");

mySocket.on('sync', function(rcvdUsers, rcvdRooms) {
  users = rcvdUsers;
  rooms = rcvdRooms;
});

/******** OBJECTS ***********/

// Note: A room name should be the same as the item name
var rooms = {
  roomName: {
    name: "name", // Identifier of the room
    services: [], // Array of allowed services
    users: {}, // Dictionary of members allowed (for easy pulling)
    bid: {
      itemName: "name",
      itemDescription: "itemDescription",
      highestBid: 0,
      bids: [],
      bidders: {},
    }
  },
}

var users = {};
var sockets = {
  uuid: "socket",
}

/********* FUNCTIONS ************/

function placeBid(roomName, bid, socket) {

  // Add/update user socket
  if (!sockets[bid.bidderID]) {
    sockets[bid.bidderID] = socket;
  }

  // Update Bid
  if(rooms.hasOwnProperty(roomName) && rooms[roomName].hasOwnProperty('services') && rooms[roomName].services.hasOwnProperty('bid')) {

    var lot = rooms[roomName].services.bid;
    var bidder = rooms[roomName].users[bid.bidderID];
    console.log("USERS", bid.bidderID);
    console.log("USERS", rooms[roomName].users);

    // Add bid to the bids array
    if (lot.bids) {
      lot.bids.push(bid);
    } else {
      lot.bids = [bid];
    }

    // Create the highest bid field if it doesn't exist
    if (!lot.highestBid) {
      lot.highestBid = 0;
    }

    // If this bid is the highest, update the highest bid & winner
    if (parseInt(bid.amount) > lot.highestBid) {
      lot.highestBid = parseInt(bid.amount);
      lot.winner = bidder;
    }

    rooms[roomName].services.bid = lot;
    console.log("Bid", bid.amount, "from", bid.bidderID);

    syncWithMain();
    updateBidders(roomName);
  } else {
    console.log("Room not found");
  }
}

function addBidToUserLog(bid) {
  var userID = bid.bidderID;
  bid.time = getDateTime();

  if (users.hasOwnProperty(userID)) {
    var user = users[userID];

    if (user.hasOwnProperty('bids')) {
      var bids = user.bids;
      bids.push(bid);
      user.bids = bids;
    } else {
      user.bids = [bid];
    }
  }

  syncWithMain();
}

// Update bidders with the most recent bid stats
function updateBidders(roomName) {

  if (rooms.hasOwnProperty(roomName) && rooms[roomName].hasOwnProperty('services') && rooms[roomName].services.hasOwnProperty('bid')) {

    var lot = rooms[roomName].services.bid;
    var bidViewers = rooms[roomName].users;

    // Send all bidders a bid update with the highest bid value and winner
    for (bidViewer in bidViewers) {
      if (sockets.hasOwnProperty(bidViewer)) {
        var socket = sockets[bidViewer];
        if (!lot.winner) {
          lot.winner = {
            'name': 'None'
          }
        }

        socket.emit('bid update', lot.highestBid, lot.winner);
      }
    }
  }
}


/****** HELPER FUNCTIONS ******/

function syncWithMain() {
 mySocket.emit('sync', users, rooms);
}

function uuid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
  }

  return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
}

function getDateTime() {

    var date = new Date();

    var hour = date.getHours();
    hour = (hour < 10 ? "0" : "") + hour;
    var min  = date.getMinutes();
    min = (min < 10 ? "0" : "") + min;
    var sec  = date.getSeconds();
    sec = (sec < 10 ? "0" : "") + sec;
    var year = date.getFullYear();
    var month = date.getMonth() + 1;
    month = (month < 10 ? "0" : "") + month;
    var day  = date.getDate();
    day = (day < 10 ? "0" : "") + day;
    return year + ":" + month + ":" + day + ":" + hour + ":" + min + ":" + sec;
}
