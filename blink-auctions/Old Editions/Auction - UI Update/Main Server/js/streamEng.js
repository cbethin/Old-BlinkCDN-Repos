var streamEng = {
  socket: null,
  serviceAddress: null
};

streamEng.setupService = function() {
  streamEng.socket = io.connect(streamEng.serviceAddress);
  console.log("Connected to StreamEng Server", streamEng.serviceAddress, roomName);

  streamEng.socket.emit('connect to stream', user.userID, roomName);
  videoInit();
}

/******** WebRTC Functionality *****/

var localVideoObject;
var remoteVideoObject;

var localStream;

var peerConnectionConfig = {
    'iceServers': [
        {'urls': 'stun:stun.services.mozilla.com'},
        {'urls': 'stun:stun.l.google.com:19302'},
    ]
};

const configOptions = {"iceServers": [{"url": "stun:stun.l.google.com:19302"},
		      {"url": "turn:35.167.210.171:3478",
				"username": "cbethin",
				"credential": "bethin"}]};

var peerConnection;

var peer1uuid = ""; // Saves UUID of Peer 1

var isCaller; // Ignore this.. it's vestigial but could be useful lol
var configuration =  null;
var peerConnectionConfig = null;

// Adding audio/video to stream
var constraints = {
  video: true,
  audio: true
}

function gotMessageFromServer(message) {
    var signal = JSON.parse(message);

    // Ignore messages from ourself
    if(signal.uuid == user.userID) return;

    if(signal.uuid == peer1uuid) {
      if(signal.type == "sdp") {
          peerConnection.setRemoteDescription(new RTCSessionDescription(signal.sdp)).then(function() {
              // Only create answers in response to offers
              if(signal.sdp.type == 'offer') {
                  console.log("Got offer")
                  peerConnection.createAnswer().then(setAndSendDescription).catch(errorHandler);
              } else {
                console.log("Got answer")
              }
          }).catch(errorHandler);
      } else if(signal.type == "ice") {
          peerConnection.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(errorHandler);
      }
    }
}

/**************** Simple Function ***********/

// Once the page has loaded, connect the JS objects to HTML objects
function videoInit() {
    localVideoObject = document.getElementById('local-video');
    remoteVideoObject = document.getElementById('remote-video');
    consoleWindow = document.getElementById('console');

    // broadcastButton = document.getElementById('broadcast1');
    hangupButton = $('#hangup');
    //startCameraButton = document.getElementById('startCamera');

    streamEng.socket.emit('join', user.userID, roomName); // Joins the server's room

    // hangupButton.disabled = true;

    window.addEventListener("beforeunload", function(e) {
        streamEng.socket.emit('disconnectServer', user.userID, roomName); // Disconnects from roomm
    }, false);

    /*************** SIGNALING *******************/

    // When it receives a ready message, send back the here message and setup the connections
    // as needed.
    streamEng.socket.on('ready', function(identifier, numClients) {
      streamEng.socket.emit('here', user.userID, roomName);

      console.log('Socket is ready');
      isCaller = identifier;

      createPeerConnection();
    });

    // When it receives a here message, save the UUID of the here message client to
    // one of the peers.
    streamEng.socket.on('here', function(new_uuid) {
      console.log("Here from " + user.userID);

      if (peer1uuid == "" && peer1uuid != new_uuid && user.userID != new_uuid) {
        peer1uuid = new_uuid;
      } else {
        console.log("Whoops");
      }
    });

    // On signal, go to gotMessageFromServer to handle the message
    streamEng.socket.on('signal', function(message) {
      console.log('Client received message:', message);
      gotMessageFromServer(message);
    });

    // Logs messages from server
    streamEng.socket.on('log', function(array) {
      console.log.apply(console, array);
    });

    streamEng.socket.on('disconnectClient', function(userID, roomName) {
      console.log(user.userID, " left the room.");
      hangup();
    });
}

// Open the local stream and create peer connection.
function setupCamera() {
  console.log("Setting up Camera");
  setupMediaStream(false);
}

// Start broadcasting to peer
function startCall() {
  hangupButton.disabled = false;
  setupMediaStream(true);
  console.log("Sending Offer");
  peerConnection.createOffer().then(setAndSendDescription).catch(errorHandler);
}

// Close connections
function hangup(userID) {
  // console.log('Ending call');
  // if (uuid == peer1uuid || !uuid ) {
  //
  // }
  //
  // //hangupButton.disabled = true;
  // //broadCastButton1.disabled = false;
  console.log("Hanging Up");
}

/********************************************/
/************* Peer Connections *************/
/********************************************/

function setupConnections(numClients) {
  // Create 1 or two connections based on # of connections to server
  if (numClients == 2) {
    console.log('Client 2 Ready.');
    createPeerConnection();
  }
}

// Get the media from camaera/microphone.
function setupMediaStream(startStream) {

  if(navigator.mediaDevices.getUserMedia) {
      navigator.mediaDevices.getUserMedia(constraints).then(getUserMediaSuccess).catch(errorHandler);
  } else {
      alert('Your browser does not support getUserMedia API');
  }

  // If you want to start the stream, addStream to connection
  if (startStream == true) {
    peerConnection.addStream(localStream);
  }
}

// Create peer connection 1
function createPeerConnection() {
  console.log("Creating Peer Connection");

  peerConnection = new RTCPeerConnection(configOptions);
  peerConnection.onicecandidate = sendIceCandidate;
  peerConnection.onaddstream = gotRemoteStream;
}


/****************************************/
/******** RTC Response Functions ********/
/****************************************/

// Set localStream object and connect local webcam to video feed
function getUserMediaSuccess(stream) {
    console.log("Success");
    localStream = stream;
    localVideoObject.src = window.URL.createObjectURL(stream);
}


function setAndSendDescription(description) {
    peerConnection.setLocalDescription(description).then(function() {
        streamEng.socket.emit('signal', JSON.stringify({'type': 'sdp', 'sdp': peerConnection.localDescription, 'uuid': user.userID}), peer1uuid, roomName);
        //serverConnection.send(JSON.stringify({'sdp': peerConnection.localDescription, 'uuid': uuid}));
    }).catch(errorHandler);
}

function sendIceCandidate(event) {
    if(event.candidate != null) {
        streamEng.socket.emit('signal', JSON.stringify({'type': 'ice', 'ice': event.candidate, 'uuid': user.userID}), peer1uuid, roomName);
        //serverConnection.send(JSON.stringify({'ice': event.candidate, 'uuid': uuid}));
    }
}

function gotRemoteStream(event) {
    console.log('Received remote stream');
    console.log(event.stream);
    remoteVideoObject.src = window.URL.createObjectURL(event.stream);
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
