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

  // $('#finalBidScreen').modal('show');

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
  $('.img-choice-image').click(function (e) {

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
    services["bid"].onTimeUpdate = onTimeUpdate;
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

function updateLot(lot) {
  // Change lot info and reset bids
  console.log("Got new lot");

  $('#item-img').attr('src', lot.itemImg)
  $('#item-name').html(function() {
    return lot.itemName;
  });

  bidArray = [];
  nextBidStr = "100";
  $('#nextBid').html(function() { return nextBidStr });

  clearGameboard();
  $('#img-box').css('background-color', WHITE);

  $('#placeBidButton').prop('disabled', false);
  $('#finalBidScreen').modal('hide');
}

function updateHTML(highestBid, winner, bidCount) {

  $('#bid-count').html(function() {
    return bidCount;
  });

  var nextBid = highestBid + 100;
  nextBidStr = nextBid.toString();

  $('#nextBid').html(function() { return nextBidStr });

  bidArray.unshift({'bid': highestBid, 'winner': winner})
  updateGameboard()

  var productBox = document.querySelector('#product-box');
  if (user.userID == winner.userID) {
    $('#img-box').css('background-color', GOLD);
    $('#placeBidButton').prop('disabled', true);
  } else {
    $('#img-box').css('background-color', WHITE);
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
      if (bidArray[i].winner.userID == user.userID) {
        return "<center><b>$" + bidArray[i].bid + "</b><br><img src=\'" + bidArray[i].winner.userImg + "\' class=userImg></center>";
      } else {
        return "<center><b>$" + bidArray[i].bid + "</b><br><img src=\'" + bidArray[i].winner.userImg + "\' class=userImg></center>";
      }
    });
  }

  for (var i=0; i < numPlayersOnScoreboard; i++) {

    if (i == 0 && bidArray[i].winner.userID == user.userID)  {
      $(gameboardLabels[i]).css('background-color', GOLD)
    } else if (bidArray[i].winner.userID == user.userID) {
      $(gameboardLabels[i]).css('background-color', GREY);
    } else {
      $(gameboardLabels[i]).css('background-color', WHITE);
    }
  }
}

function clearGameboard() {
  // Reset background color of product box
  $('#product-box').css('background-color', WHITE);

  // Reset the game row
  for (var i = 0; i < 6; i++) {

    $(gameboardLabels[i]).html(function() { return "" });
    $(gameboardLabels[i]).css('background-color', WHITE);
  }

  // Reset the bid count & timer
  $('#bid-count').html(function() { return "0"});
  $('#timer').html(function() {return "2:00"});
  $('#nextBid')


}

function onFinalBidUpdate(highestBid, winner, bidCount, lot) {
  console.log("Final Bid Update");

  var winYesNo;
  var textColor = BLACK;

  if (winner.userID == user.userID) {
    winYesNo = "You Win";
    textColor = GOLD;
  } else {
    winYesNo = "You Lost";
  };

  console.log(lot)
  $('#timer').html(function() { return "0:00" });
  $('#bid-count').html(function() { return "0" });
  $('#nextBid').html(function() { return "100" });
  $('#nextImg').attr('src', lot.itemImg);
  $('#nextName').html(function() { return lot.itemName });
  $('#nextDescription').html(function() { return lot.description });
  $('#winner').html(function() { return winYesNo });
  $('#winner').css('color', textColor);

  nextBidStr = "100";
}

function onTimeUpdate(timeLeft, isOnBreak) {

  if (isOnBreak) {
    $('#finalBidScreen').modal('show');
    $('#break-countdown').html(function() {
      return convertToMinutesAndSeconds(timeLeft);
    });
  } else {
    $('#finalBidScreen').modal('hide');
    $('#timer').html(function() {
      return convertToMinutesAndSeconds(timeLeft);
    });

    if (timeLeft < 30) {
      $('#timer').css('color', WARNING_RED);
    } else {
      $('#timer').css('color', BLACK);
    }
  }


}

function convertToMinutesAndSeconds(timeInSeconds) {
  var minutes = Math.floor(timeInSeconds / 60);
  var seconds = timeInSeconds - minutes * 60;

  var secondString = "";

  if (seconds < 10) {
    secondString = "0" + seconds.toString();
  } else {
    secondString = seconds.toString();
  }

  return minutes.toString() + ":" + secondString;
}

function setupVideo() {
  console.log("SETUP VIDEO");
  // videoInit();
}


var WARNING_RED = "#E82A42";
var GOLD = "#FFE0B3";
var GREY = "#F0F0F0";
var WHITE = "#FFFFFF";
var BLACK = "#000000";
