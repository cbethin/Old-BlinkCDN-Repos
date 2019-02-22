const HTTPS_PORT = 5000;

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

      socket.on('disconnectServer', function(uuid, roomName) {
        onDisconnect(uuid, roomName);
      });

      socket.on('connect to stream', function(userID, roomName) {
        onJoin(userID, socket, roomName);
      });
});

/******* SETUP MAIN SERVER CONNECTION *********/

var io_client = require('socket.io-client');
var mySocket = io_client.connect("http://localhost:3001");
mySocket.emit('connect service', "https://localhost:5000", "stream");

mySocket.on('sync', function(rcvdUsers, rcvdRooms) {
  users = rcvdUsers;
  rooms = rcvdRooms;
});

/******* FUNCTIONALITY **********/

function onSignal(message, destUuid, roomName, socket) {
  var signal = JSON.parse(message);
  var room = streamRooms[roomName];

  for (var i = 0; i < room.clients.length; i++) {
    if (room.clients[i].uuid == destUuid) {
      console.log("Sending", signal.type, " from ", socket.id, " to ", room.clients[i].uuid)
      room.clients[i].socket.emit('signal', message, socket.id);
    };
  };
}

function onDisconnect(uuid, roomName) {
    console.log(uuid, "Disconnecting");

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
           if (clientsInRoom[i].uuid == uuid) {
              // If this is the client, just remove them from the room
              clientsInRoom.splice(i, 1);
              streamRooms[roomName].clients = clientsInRoom;
           } else {
              // If this isn't the client, let them know the other client is leaving
              clientsInRoom[i].socket.emit('disconnectClient', uuid, roomName);
              console.log("Sent disconnect")
           }
        }
      }
     }
}

function onJoin(uuid, socket, roomName) {

   if (!streamRooms[roomName]) {
    //If the room does not exist, create it
    console.log(socket.id, " created new room with id:", roomName);
    streamRooms[roomName] = {
      clients: [{"uuid": uuid, "socket": socket}]
    }
  } else {
    // If rooms exist, and the most recent room only has one client,
    // add this client to the room
    clientsInThisRoom = streamRooms[roomName].clients
    clientsInThisRoom.push({'uuid': uuid, 'socket': socket});
    streamRooms[roomName].clients = clientsInThisRoom;

    // SEND A 'Here' to the stream server <-- !! FIND A BETTER WAY TO INDENTIFY !! -->
    if (clientsInThisRoom[0].uuid != uuid) {
      clientsInThisRoom[0].socket.emit('here', uuid);
    }

    // Let everyone know the system is ready
    for (var i = 0; i < clientsInThisRoom.length; i++) {
      clientsInThisRoom[i].socket.emit('ready');
    }

    console.log(socket.id, " joined the room ", roomName);
  }
}
