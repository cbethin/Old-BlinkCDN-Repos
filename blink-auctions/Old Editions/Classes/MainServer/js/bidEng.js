// BidEng Object
var bidEng = {
  // Properties
  socket: null,
  serviceAddress: null,
  // Functions
  onBidUpdate: null,
}

bidEng.setupService = function() {
  bidEng.socket = io.connect(bidEng.serviceAddress);
  console.log("Connected to BidEng Server", bidEng.serviceAddress);

  bidEng.socket.on('bid update', function(highestBid, winner) {
    console.log("got update");
    bidEng.onBidUpdate(highestBid, winner);
  });

  bidEng.socket.emit('connect to bid', user.userID, roomName);
}

bidEng.placeBid = function(bidValue) {
  console.log("Bidding:", bidValue);

  var bid = {
    roomName: roomName,
    bidderID: user.userID,
    amount: bidValue
  };

  bidEng.socket.emit('place bid', bid);
}
