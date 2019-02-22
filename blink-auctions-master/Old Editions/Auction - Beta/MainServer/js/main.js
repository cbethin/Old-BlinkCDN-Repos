// Connects to socket.io server
var socket;
var uuid;
var roomName = 'helloAdele';
console.log("Connected");

// Setup HTML Objects
var button;
var bidValueInput;

var user = {};

var services = {
  "bid": bidEng,
  "stream": streamViewEng
}

$(document).ready(function() {

  // Setup Socket
  setupSocket();

  // Setup HTML Objects
  button = $("button");
  bidValueInput = $("#placeBidInput");

  user.name = prompt("Please enter your name", "full name");
  user.email = "cbethin@stevens.edu";

  // Join main SERVER
  socket.emit('create user', user, roomName);

  // Place a bid
  button.click(function() {
    var bidValue = bidValueInput.val();
    bidEng.placeBid(bidValue);
  });

});

/******* SOCKET ********/

function setupSocket() {

  socket = io.connect();

  socket.on('created user', function(userID) {
    user.userID = userID;

    // Send Join Bid System Message
    services["bid"].onBidUpdate = updateHTML;
    socket.emit('join service', user.userID, 'bid', roomName);

    // Send join stream system Message
    socket.emit('join service', user.userID, 'stream', roomName);
  });

  socket.on('joined service', function(userID, serviceType, serviceAddress) {
    var engine = services[serviceType];
    engine.serviceAddress = serviceAddress;

    engine.setupService();
  });

  socket.on('bid update', function() {
    console.log("UPDATE");
  });

  console.log("Setup socket");
}

/****** FUNCTIONS ********/

var lastBid;

function updateHTML(highestBid, winner) {
  console.log(highestBid, winner);

  if (lastBid != highestBid) {
    shiftWinners(highestBid, winner.name);
  };


  $('#highestBid').html(function() {
    return "Highest Bid: " + highestBid;
  });

  $('#bidWinner').html(function() {
    return "Highest Bidder: " + winner.name;
  });

  lastBid = highestBid;
  lastWinner = winner;
}

function shiftWinners(highestBid, winnerName) {

  $('#sixth-place').html(function() {
    return $("#fifth-place").html();
  });

  $('#fifth-place').html(function() {
    return $("#fourth-place").html();
  });

  $('#fourth-place').html(function() {
    return $("#third-place").html();
  });

  $('#third-place').html(function() {
    return $("#second-place").html();
  });

  $('#second-place').html(function() {
    return $("#first-place").html();
  });

  $('#first-place').html(function() {
    return "<h5>$ " + highestBid + "</h5><h5>" + winnerName + "</h5>";
  });
}

function setupVideo() {
  console.log("SETUP VIDEO");
  // videoInit();
}
