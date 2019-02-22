// Connects to socket.io server
var socket;
var uuid;
var roomName = 'helloAdele';
console.log("Connected");

// Setup HTML Objects
var button;
var nextBidStr;

/* user = {
  name:
  userImg:
  userID:
}*/
var user = {};

var services = {
  "bid": bidEng,
  "stream": streamViewEng
}

// Bid Array
var bidArray = [];
const gameboardLabels = ["#first", "#second", "#third", "#fourth", "#fifth", "#sixth"];

$(document).ready(function() {

  // Setup Socket
  setupSocket();

  // Setup HTML Objects
  nextBidStr = $("#nextBid").text();

  // When the placeBidButton is pressed, place a bid
  $("#placeBidButton").click(function() {
    var bidValue = parseInt(nextBidStr);
    console.log(nextBidStr);
    bidEng.placeBid(bidValue);
  });

  // Show the User Setup Window
  $('#userSetupScreen').modal('show');

  // When an image is clicked, if it is selected then deselect it. If it's deselected then select it
  // and deslect all the others. Update the userImg.
  $('.imgSelect').click(function (e) {

    var target = $(e.target), article;
    console.log("Target:", target);

    $('.clicked-img').each(function() {
      if ($(this).attr('src') != target.attr('src')) {
        $(this).removeClass('clicked-img');
        console.log($(this), target);
      }
    });


    if (target.hasClass('clicked-img')) {
      target.removeClass('clicked-img');
    } else {
      target.addClass('clicked-img');
      user.userImg = target.attr('src');
    }
  });

  // When User Input Window is dismissed, grab user info and join the main SERVER
  // by emitting a 'create user' message
  $('#modal-done-button').click(function() {

    user.name = $('#usernameInput').val();
    $('#userSetupScreen').modal('hide');

    socket.emit('create user', user, roomName);
  });

});

/******* SOCKET ********/

function setupSocket() {

  socket = io.connect();

  socket.on('created user', function(userID) {
    user.userID = userID;

    // Send Join Bid System Message
    services["bid"].onBidUpdate = updateHTML;
    services["bid"].onNewLot = updateLot;
    services["bid"].onFinalBidUpdate = onFinalBidUpdate;
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

  $('#userSetupDoneButton').click(function() {
    tryToSetupUser();
  });
}

/****** FUNCTIONS ********/

function tryToSetupUser() {
  // Is the name input empty?
  // Is an image selected?
  // Get that data and store it
}

function updateLot(lot) {
  // Change lot info and reset bids
  $('#item-img').attr('src', lot.itemImg)

  $('#item-name').html(function() {
    return lot.itemName;
  });

  bidArray = [];
  nextBidStr = "100";
  $('#nextBid').html(function() { return nextBidStr });

  clearGameboard();

  $('#placeBidButton').prop('disabled', false);
}

function updateHTML(highestBid, winner, bidCount) {

  $('#bid-count').html(function() {
    return bidCount;
  });

  var nextBid = highestBid + 100;
  nextBidStr = nextBid.toString();

  $('#nextBid').html(function() { return nextBidStr });

  bidArray.unshift({'bid': highestBid, 'winner': winner})
  console.log("Bid Array:");
  updateGameboard()

  var productBox = document.querySelector('#product-box');
  if (user.userID == winner.userID) {
    productBox.style.backgroundColor = "#FFE0B3";
    $('#placeBidButton').prop('disabled', true);
  } else {
    productBox.style.backgroundColor = "#FFFFFF";
    $('#placeBidButton').prop('disabled', false);
  }

}

function updateGameboard() {
  var numPlayersOnScoreboard = 0;

  // Make sure the scoreboard only tries to put at most 6 players on it
  if (bidArray.length > 6) {
    numPlayersOnScoreboard = 6;
  } else {
    numPlayersOnScoreboard = bidArray.length;
  }

  for (var i = 0; i < numPlayersOnScoreboard; i++) {
    $(gameboardLabels[i]).html(function() {
      return "<center><b>$" + bidArray[i].bid + "</b><br><img src=\'" + bidArray[i].winner.userImg + "\' class=userImg></center>";
    });
  }

  for (var i=0; i < numPlayersOnScoreboard; i++) {

    if (i == 0 && bidArray[i].winner.userID == user.userID)  {
      $(gameboardLabels[i]).css('background-color', "#FFE0B3")
    } else if (bidArray[i].winner.userID == user.userID) {
      $(gameboardLabels[i]).css('background-color', "#F0F0F0");
    } else {
      $(gameboardLabels[i]).css('background-color', "#FFFFFF");
    }
  }
}

function clearGameboard() {
  // Reset background color of product box
  $('#product-box').css('background-color', '#FFFFFF');

  // Reset the game row
  for (var i = 0; i < 6; i++) {

    $(gameboardLabels[i]).html(function() { return "" });
    $(gameboardLabels[i]).css('background-color', '#FFFFFF');
  }

  // Reset the bid count & timer
  $('#bid-count').html(function() { return "0"});
  $('#timer').html(function() {return "2:00"});
  $('#nextBid')


}

function onFinalBidUpdate() {
  console.log("Final Bid Update");

  $('#timer').html(function() { return "0:00" });
  $('#bid-count').html(function() { return "0" });
  $('#nextBid').html(function() { return "100" });
  nextBidStr = "100";

}


function setupVideo() {
  console.log("SETUP VIDEO");
  // videoInit();
}
