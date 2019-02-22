//STATUS: WOrking
var localVideoObject;
var remoteVideoObject;
var broadcastButton;

var isBroadcaster = false;
var broadcaster;
var roomName = "helloAdele";
var localStreams = {};
var peers = {};
var sendToPeerValue = -1;

const configOptions = {"iceServers": [{"url": "stun:stun.l.google.com:19302"},
              { url: 'turn:numb.viagenie.ca',
                credential: 'enter1234',
                username: 'bethin.charles@yahoo.com'
              }]};

var peers = [];

var constraints = {
  video: true,
  audio: true
}

function gotMessageFromServer(message) {
    var signal = message;
    var peerNumber = -1;

    // Ignore messages from ourself
    if(signal.userID == user.userID) {
      console.log("Received from self");
      return;
    }


    if (isBroadcaster) {
      for (var i=0; i < peers.length; i++) {
        if (peers[i].userID == signal.userID) {
          peerNumber = i;
          break;
        }
      }

      if (peers[peerNumber].userID == signal.userID) {
        //sendToPeerValue = peer.number;

        if(signal.type == "sdp") {
            peers[peerNumber].peerConnection.setRemoteDescription(new RTCSessionDescription(signal.sdp)).then(function() {
                // Only create answers in response to offers
                if(signal.sdp.type == 'offer') {
                    console.log("Got offer")
                    sendToPeerValue = peerNumber;
                    peers[peerNumber].peerConnection.createAnswer().then(setAndSendDescription).catch(errorHandler);
                } else {
                  console.log("Got answer")
                }
            }).catch(errorHandler);
        } else if(signal.type == "ice") {
            peers[peerNumber].peerConnection.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(errorHandler);
            console.log("Signal Ice:", signal.ice);
        }
      }
    } else {
      if (broadcaster.castID == signal.userID) {

        if(signal.type == "sdp") {

            try {
              broadcaster.peerConnection.setRemoteDescription(new RTCSessionDescription(signal.sdp)).then(function() {
                  // Only create answers in response to offers
                  if(signal.sdp.type == 'offer') {
                      console.log("Got offer")
                      sendToPeerValue = -10;
                      broadcaster.peerConnection.createAnswer().then(setAndSendDescription).catch(errorHandler);
                  } else {
                    console.log("Got answer")
                  }
              }).catch(errorHandler);
            }
            catch(err) {
              console.log(err);
            }

        } else if(signal.type == "ice") {
            broadcaster.peerConnection.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(errorHandler);
            console.log("Signal Ice:", signal.ice);
        }
      }
    }



}

function joinRoom() {
  // It only runs two of each cuz of that error;
  try {
    startCall();
  } catch(err) {
    console.log("Error:", err)
  }
  try {
    startCall();
  } catch(err) {
    console.log("Error:", err)
  }
}

function startCall() {
  // for (var i=0; i<peers.length; i++) {
    setupMediaStream(true, peers.length-1);
  // }

}


// Get the media from camaera/microphone.
function setupMediaStream(startStream, peerNumber) {

  if(navigator.mediaDevices.getUserMedia) {
      navigator.mediaDevices.getUserMedia(constraints).then(function(stream) {
        localStreams[peerNumber] = stream;
        localVideoObject.src = window.URL.createObjectURL(stream);

        // If you want to start the stream, addStream to connection
        if (startStream == true) {
            console.log("Adding media stream to:", peerNumber);
            peers[peerNumber].peerConnection.addStream(localStreams[peerNumber]);

            console.log("Sending Offer");
            sendToPeerValue = peerNumber;
            peers[peerNumber].peerConnection.createOffer().then(setAndSendDescription).catch(errorHandler);
        }
      }).catch(errorHandler);
  } else {
      alert('Your browser does not support getUserMedia API');
  }
}

// Create peer connection 1
function createPeerConnection(peerUserID, peerNumber) {
  console.log("Creating Peer Connection");

  var newPeerConnection = new RTCPeerConnection(configOptions);
  newPeerConnection.onicecandidate = function(event) {
    if(event.candidate != null) {
        streamEng.socket.emit('signal', {'type': 'ice', 'ice': event.candidate, 'userID': user.userID}, peerUserID, roomName);
    }
  };

  newPeerConnection.onaddstream = function(event) {
    console.log('Received remote stream');
    console.log(event.stream);
    remoteVideoObject.src = window.URL.createObjectURL(event.stream);
  };

  return newPeerConnection;
}

function setAndSendDescription(description) {

  if (sendToPeerValue == -10) {
    broadcaster.peerConnection.setLocalDescription(description).then(function() {
      console.log("Local Description:", broadcaster.peerConnection.localDescription);
        streamEng.socket.emit('signal', {'type': 'sdp', 'sdp': broadcaster.peerConnection.localDescription, 'userID': user.userID}, broadcaster.castID, roomName);
    }).catch(errorHandler);
  } else {
    peers[sendToPeerValue].peerConnection.setLocalDescription(description).then(function() {
        console.log("Local Description:", peers[sendToPeerValue].peerConnection.localDescription);
        streamEng.socket.emit('signal', {'type': 'sdp', 'sdp': peers[sendToPeerValue].peerConnection.localDescription, 'userID': user.userID}, peers[sendToPeerValue].userID, roomName);
        console.log("Sent local description to", peers[sendToPeerValue].userID);
    }).catch(errorHandler);
  }
}


// StreamCast Eng Stuff

var streamEng = {
  socket: null,
  serviceAddress: null
}

streamEng.setupService = function() {
  setupPage();
  streamEng.socket = io.connect(streamEng.serviceAddress);
  console.log("Connected to Stream Server", streamEng.serviceAddress, roomName);

  streamEng.socket.emit('connect to stream', user.userID, roomName, isBroadcaster);

  // When it receives a here message, add user to peers
  streamEng.socket.on('here', function(clientID) {
    console.log("Here from " + clientID);

    //if clientID is still blank, AND if clientID doesn't exist yet AND this device isn't the userID
    if (!peers[clientID] && user.userID != clientID) {
      var newPeerConnection = createPeerConnection(clientID, peers.length);
      peers.push({
        "userID": clientID,
        "number": (peers.length),
        "peerConnection": newPeerConnection
      });
    } else {
      console.log("Whoops");
    }

    joinRoom();
  });

  // The broadcaster is ready to stream, create a PC for it
  streamEng.socket.on('ready', function(castID) {
    console.log("Broadcaster is ready.");

    var newPeerConnection = createPeerConnection(castID);
    broadcaster = {
      "castID": castID,
      "peerConnection": newPeerConnection
    };
  });

  // On signal, go to gotMessageFromServer to handle the message
  streamEng.socket.on('signal', function(message) {
    console.log('Client received message:', message);
    gotMessageFromServer(message);
  });

  // Handle client disconnect
  streamEng.socket.on('disconnect client', function(userID, roomName) {
    console.log(user.userID, " left the room.");
  });

}

// Setup DOM elements and responses
function setupPage() {
    localVideoObject = document.getElementById('local-video');
    remoteVideoObject = document.getElementById('remote-video');


    // If client is going to disconnect, let server know
    window.addEventListener("beforeunload", function(e) {
        streamEng.socket.emit('disconnect client', user.userID, roomName); // Disconnects from roomm
    }, false);
}

///////////////////
function errorHandler(error) {
    console.log(error);
}
