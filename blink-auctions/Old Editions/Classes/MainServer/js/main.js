// Connects to socket.io server
var socket;
var uuid;
var roomName = 'CPE360';
console.log("Connected");

// Setup HTML Objects
var button;
var bidValueInput;

/* user = {
  name:
  email:
  identifier:
}*/
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

  user.name = "unknown"
  user.email = "unknown";

  // Join main SERVER
  socket.emit('create user', user, roomName);

  // Send join Stream system message
  // socket.emit('join system', user, 'stream', roomName);
  // NEED TO REWRITE SO THAT THE SERVER DOES NOT GENERATE NEW UUID FOR EVER SERVICE

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

function updateHTML(highestBid, winner) {
  console.log(highestBid, winner);

  $('#highestBid').html(function() {
    return "Highest Bid: " + highestBid;
  });

  $('#bidWinner').html(function() {
    return "Highest Bidder: " + winner.name;
  });
}

function setupVideo() {
  console.log("SETUP VIDEO");
  // videoInit();
}
