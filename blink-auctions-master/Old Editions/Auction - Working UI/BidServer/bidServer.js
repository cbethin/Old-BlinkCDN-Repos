const HTTPS_PORT = 4000;
const BREAK_TIME = 60000; // 60 seconds
const AUCTION_TIME = 120000; // 2 minutes

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
    setupBidUser(userID, roomName);
  });

  socket.on('place bid', function(bid) {
    var roomName = bid.roomName;

    placeBid(roomName, bid, socket);
    addBidToUserLog(bid);
  });

});

/******* SETUP MAIN SERVER CONNECTION *********/

var io_client = require('socket.io-client');
var mySocket = io_client.connect("http://localhost:3001");
mySocket.emit('connect service', "https://localhost:4000", "bid");

mySocket.on('sync', function(rcvdUsers, rcvdRooms) {
  // See what new rooms are added
  var timersToSet = [];
  for (roomName in rcvdRooms) {
    if (!rooms[roomName]) {
      timersToSet.push(roomName);
    }
  }

  // Update this server's info
  users = rcvdUsers;
  rooms = rcvdRooms;

  // Start auction timers
  timersToSet.forEach(function(roomName) {
    console.log("Starting", roomName, "auction timer");
    startAuctionTimer(roomName);
  });

});

/******** OBJECTS ***********/

var rooms = {
  /* roomName: {
    name: "name", // Identifier of the room
    services: [], // Array of allowed services
    users: {}, // Dictionary of members allowed (for easy pulling)
    bid: {
      lots: [{
        itemName: "Mona Lisa",
        description: "Does it need a description?",
        itemImg: "img/items/monalisa.png"
        },
        {
        itemName: "Rolex Daytona 116500",
        description: "High end fancy watch.",
        itemImg: "img/items/rolex.png"
        }],
      currentLotNumber: 0,
      highestBid: 0,
      bidCount: 0,
      bids: [],
    }
  }, */
}

var auctionTimers = {
  // roomName: {
  //   auctionTimer: null,
  //   breakTimer: null,
  // },
}

var users = {};
var sockets = {
  uuid: "socket",
}

var mostRecentBids = [];

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
      lot.bidCount++;
      lot.highestBid = parseInt(bid.amount);
      lot.winner = bidder;
      mostRecentBids.unshift({'bid': lot.highestBid, 'winner': lot.winner});
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

        socket.emit('bid update', lot.highestBid, lot.winner, lot.bidCount);
      }
    }
  }
}

// Send info needed to setup bid to user
function setupBidUser(userID, roomName) {

  var topSixBids = []
  if (sockets.hasOwnProperty(userID)) {

    var socket = sockets[userID]
    var numBidsToSend = 0;

    if (mostRecentBids.length > 5) {
      numBidsToSend = 5;
    } else {
      numBidsToSend = mostRecentBids.length;
    }

    for(var i = 0; i <= numBidsToSend; i++) {
      if (mostRecentBids[i]) {
        topSixBids.push(mostRecentBids[i]);
      }
    }

    var currentLotNumber = rooms[roomName].services.bid.currentLotNumber;
    var lots = rooms[roomName].services.bid.lots
    socket.emit('new lot', lots[currentLotNumber%lots.length])
    socket.emit('bid setup', mostRecentBids, rooms[roomName].services.bid.bidCount);
  }

}

/******* TIMER FUNCTIONALITY **********/

function startAuctionTimer(roomName) {

  if (!auctionTimers[roomName]) {
    auctionTimers[roomName] = {
      auctionTimer: null,
      breakTimer: null
    }
  }

  /*auctionTimers[roomName].auctionTimer =*/ setTimeout(function() {
    console.log("Auction Timer Stopped");
    onEndOfAuction(roomName);
  }, AUCTION_TIME);

}

function onEndOfAuction(roomName) {

  for (bidViewer in rooms[roomName].users) {
    if (sockets.hasOwnProperty(bidViewer)) {
      var socket = sockets[bidViewer];
      var bidRoom = rooms[roomName].services.bid;
      if (!bidRoom.winner) {
        bidRoom.winner = {
          'name': 'None'
        }
      }

      socket.emit('final bid', bidRoom.highestBid, bidRoom.winner, bidRoom.bidCount);
    }
  }

  rooms[roomName].services.bid.bids = [];
  rooms[roomName].services.bid.bidCount = 0;
  rooms[roomName].services.bid.highestBid = 0;

  startBreakTimer(roomName);
}

function startBreakTimer(roomName) {
  /*auctionTimers[roomName].breakTimer =*/ setTimeout(function() {
    console.log("Break timer stopped");
    onEndOfBreak(roomName);
  }, BREAK_TIME);
}

function onEndOfBreak(roomName) {

  rooms[roomName].services.bid.currentLotNumber++;
  var currentLotNumber = rooms[roomName].services.bid.currentLotNumber
  var bidRoom = rooms[roomName].services.bid;

  for (bidViewer in rooms[roomName].users) {
    if (sockets.hasOwnProperty(bidViewer)) {
      var socket = sockets[bidViewer];
      if (!bidRoom.winner) {
        bidRoom.winner = {
          'name': 'None'
        }
      }

      console.log("Sending new lots");
      socket.emit('new lot', bidRoom.lots[currentLotNumber%2]);
    }
  }

  console.log("Restarting timer for:", roomName);
  startAuctionTimer(roomName);

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
