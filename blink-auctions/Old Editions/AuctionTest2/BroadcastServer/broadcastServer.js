const HTTPS_PORT = 443;

const nodeStatic = require('node-static');
const https = require('https');
const socketIO = require('socket.io');
const fs = require('fs');
const os = require('os');

/******** OBJECTS ***********/

// Rooms
var streamRooms = {};

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

      socket.on('here', function(userID, roomName) {
        console.log("Here from: ", userID);
        for (var i = 0; i < streamRooms[roomName].clients.length; i++) {
          streamRooms[roomName].clients[i].socket.emit('here', userID);
          console.log("Send \'here\' to ", streamRooms[roomName].clients[i].socket.id);
        }
      })

      socket.on('signal', function(message, destUuid, roomName) {
        onSignal(message, destUuid, roomName, socket);
      });

      socket.on('disconnection', function() {
        console.log(socket.id, ' disconnected!')
      });

      socket.on('disconnectServer', function(userID, roomName) {
        onDisconnect(userID, roomName);
      });

      socket.on('connect to stream', function(userID, roomName, isBroadcaster) {
        onJoin(userID, socket, roomName, isBroadcaster);
      });
});

/******* SETUP MAIN SERVER CONNECTION *********/

var io_client = require('socket.io-client');
var mySocket = io_client.connect("http://bid.blinkcdn.com");
mySocket.emit('connect service', "https://streamserver.blinkcdn.com", "stream");

mySocket.on('sync', function(rcvdUsers, rcvdRooms) {
  users = rcvdUsers;
  rooms = rcvdRooms;
});

/******* FUNCTIONALITY **********/

function onSignal(message, destUuid, roomName, socket) {
  var signal = JSON.parse(message);
  var room = streamRooms[roomName];

  if (destUuid == room.castID) {
    console.log("Sending", signal.type, "to broadcaster.");
    room.castSocket.emit('signal', message);
    return;
  }

  for (var i = 0; i < room.clients.length; i++) {
    if (room.clients[i].userID == destUuid) {
      console.log("Sending", signal.type, " from ", socket.id, " to ", room.clients[i].userID)
      room.clients[i].socket.emit('signal', message);
    };
  };
}

function onDisconnect(userID, roomName) {
    console.log(userID, "Disconnecting");

    if(streamRooms[roomName]) {
      var clientsInRoom = streamRooms[roomName].clients

      if(clientsInRoom.length == 1) {
        streamRooms[roomName] = null;
        return;
      }

      if (clientsInRoom.length == 1) {
        streamRooms[roomName] = null;
        return;
      } else {
        for(var i = 0; i < clientsInRoom.length; i++) {
           if (clientsInRoom[i].userID == userID) {
              // If this is the client, just remove them from the room
              clientsInRoom.splice(i, 1);
              streamRooms[roomName].clients = clientsInRoom;
           } else {
              // If this isn't the client, let them know the other client is leaving
              clientsInRoom[i].socket.emit('disconnectClient', userID, roomName);
              console.log("Sent disconnect")
           }
        }
      }
     }
}

function onJoin(userID, socket, roomName, isBroadcaster) {

  // IF it is a broadcaster, setup as the broadcaster;
  if (isBroadcaster == true) {

    if (!streamRooms[roomName]) {
      streamRooms[roomName] = {
        clients: [{"userID": userID, "socket": socket}],
        castSocket: null,
        castID: null
      }
    }

    streamRooms[roomName].castSocket = socket;
    streamRooms[roomName].castID = userID;

    // Send a message to clients letting them know broadcaster joined
    // Send message to broadcaster letting them know clients are there
    for (var i = 0; i < streamRooms[roomName].clients.length; i++) {
      streamRooms[roomName].clients[i].socket.emit('ready', streamRooms[roomName].castID);
      streamRooms[roomName].castSocket.emit('here', streamRooms[roomName].clients[i].userID)
    }

    console.log("Streamer joined the session:", roomName);
    return;
  }

  // If not, then if the room doesn't exist, create the room
  else if (!streamRooms[roomName]) {
    console.log("Client created room:", roomName);
    streamRooms[roomName] = {
      clients: [{"userID": userID, "socket": socket}]
    }
  }

  // If room exists, then add the client to the room
  else {
    streamRooms[roomName].clients.push({'userID': userID, 'socket': socket});

    // Let the broadcaster know someone else joined and let client know the system is ready
    if (streamRooms[roomName].castID != null) {
      streamRooms[roomName].castSocket.emit('here', userID);
      socket.emit('ready', streamRooms[roomName].castID);
    }

    // If braodcaster doesn't exist, wait for them to join.

    console.log(socket.id, " joined the room ", roomName);
  }

}
