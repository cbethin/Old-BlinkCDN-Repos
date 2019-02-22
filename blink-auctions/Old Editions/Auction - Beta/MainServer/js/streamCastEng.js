var streamCastEng = {
  socket: null,
  serviceAddress: null
};

/******** WebRTC Functionality *****/

var localVideoObject;
var localStreams = {};
var broadcastButton;
var hangupButton;

var peerConnectionConfig = {
    'iceServers': [
        {'urls': 'stun:stun.services.mozilla.com'},
        {'urls': 'stun:stun.l.google.com:19302'},
    ]
};

// const configOptions = {"iceServers": [{"url": "stun:stun.l.google.com:19302"},
// 		      {"url": "turn:35.167.210.171:3478",
// 				"username": "cbethin",
// 				"credential": "bethin"}]};

var peers = {};
var sendToPeerValue;

// Adding audio/video to stream
var constraints = {
  video: true,
  audio: true
}

function gotMessageFromServer(message) {
    var signal = JSON.parse(message);
    var peerUserID = signal.uuid;

    // Ignore messages from ourself
    if(signal.uuid == user.userID) return;

    if(signal.type == "sdp") {
        peers[peerUserID].peerConnection.setRemoteDescription(new RTCSessionDescription(signal.sdp)).then(function() {
            // Only create answers in response to offers
            if(signal.sdp.type == 'offer') {
                console.log("Got offer")
                sendToPeerValue = peerUserID;
                peers[peerUserID].peerConnection.createAnswer().then(setAndSendDescription).catch(errorHandler);
            } else {
              console.log("Got answer");
            }
        }).catch(errorHandler);
    } else if(signal.type == "ice") {
        peers[peerUserID].peerConnection.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(errorHandler);
        console.log("Signal Ice:", signal.ice);
    }

}

/**************** Simple Function ***********/

streamCastEng.setupService = function() {
  streamCastEng.socket = io.connect(streamCastEng.serviceAddress);
  console.log("Connected to Stream Server", streamCastEng.serviceAddress, roomName);

  streamCastEng.socket.emit('connect to stream', user.userID, roomName, true);

  // When it receives a here message, add user to peers
  streamCastEng.socket.on('here', function(clientID) {
    console.log("Here from " + clientID);

    //If the peer exists, and clientID is not this device's userID, and there are less than 2 peers
    if (!peers.hasOwnProperty(clientID) && user.userID != clientID /*&& Object.keys(peers).length < 2*/) {
      var newPeerConnection = createPeerConnection(clientID);
      peers[clientID] = {
        "userID": clientID,
        "peerConnection": newPeerConnection
      };

      startBroadcastTo(clientID);
    }
  });

  // On signal, go to gotMessageFromServer to handle the message
  streamCastEng.socket.on('signal', function(message) {
    console.log('Client received message:', message);
    gotMessageFromServer(message);
  });

  // Handle client disconnect
  streamCastEng.socket.on('disconnect client', function(userID, roomName) {
    console.log(user.userID, " left the room.");
  });

  setupPage();
}

streamCastEng.endService = function() {
  console.log("Ending service");
  remoteVideoObject.src = null;
  peers = {};
  sendToPeerValue = null;
}

streamCastEng.disconnect = function() {
  streamCastEng.socket.emit('cast disconnect', user.userID, roomName);
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

// Start Broadcasting to a given peer
function startBroadcastTo(peerID) {
  setupMediaStream(true, peerID);
  // sendToPeerValue = peerID;
  // peers[peerID].peerConnection.createOffer().then(setAndSendDescription).catch(errorHandler);
}

// Close connections
function hangup(userID) {
  // console.log('Ending call');
  // if (uuid == peer1uuid || !uuid ) {
  //
  // }
  // //hangupButton.disabled = true;
  // //broadCastButton1.disabled = false;
  console.log("Hanging Up");
}

/********************************************/
/************* Peer Connections *************/
/********************************************/

// Get the media from camaera/microphone.
function setupMediaStream(startStream, peerID) {

  if(navigator.mediaDevices.getUserMedia) {

      navigator.mediaDevices.getUserMedia(constraints).then(function(stream) {
        localStreams[peerID] = stream;
        localVideoObject.src = window.URL.createObjectURL(stream);

        peers[peerID].peerConnection.addStream(localStreams[peerID]);
        console.log("Adding media stream to:", peerID);

        // If you want to start the stream, addStream to connection
        if (startStream == true) {
          console.log("Sending offer.");
          sendToPeerValue = peerID;
          peers[peerID].peerConnection.createOffer().then(setAndSendDescription).catch(errorHandler);
        }

      }).catch(errorHandler);
  } else {
      alert('Your browser does not support getUserMedia API');
  }
}

// Create peer connection 1
function createPeerConnection(peerUserID) {
  console.log("Creating Peer Connection");

  var newPeerConnection;
  newPeerConnection = new RTCPeerConnection(peerConnectionConfig);
  newPeerConnection.onicecandidate = function(event) {
    if(event.candidate != null) {
        socket.emit('signal', JSON.stringify({'type': 'ice', 'ice': event.candidate, 'uuid': user.userID}), peerUserID, roomName);
    }
  };

  newPeerConnection.onaddstream = function(event) {
    console.log('Received remote stream');
    console.log(event.stream);
    remoteVideoObjects[peerNumber].src = window.URL.createObjectURL(event.stream);
  };

  console.log("Created Object:", newPeerConnection);
  return newPeerConnection;
}


/****************************************/
/******** RTC Response Functions ********/
/****************************************/


function setAndSendDescription(description) {

  peers[sendToPeerValue].peerConnection.setLocalDescription(description).then(function() {
      streamCastEng.socket.emit('signal', JSON.stringify({'type': 'sdp', 'sdp': peers[sendToPeerValue].peerConnection.localDescription, 'uuid': user.userID}), peers[sendToPeerValue].userID, roomName);
      console.log("Send to Peer", sendToPeerValue);
  }).catch(errorHandler);

}

/****************
    Helper Functions
****************/


function errorHandler(error) {
    console.log(error);
}

function randomToken() {
  return Math.floor((1 + Math.random()) * 1e16).toString(16).substring(1);
}

// Taken from http://stackoverflow.com/a/105074/515584
// Strictly speaking, it's not a real UUID, but it gets the job done here
function uuid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
  }

  return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
}
