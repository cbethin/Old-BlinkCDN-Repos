const HTTPS_PORT = 4000;
const BREAK_TIME = 15000; // 60 seconds
const AUCTION_TIME = 12000; // 2 minutes

/*******Start*********/

const nodeStatic = require('node-static');
const https = require('https');
const http = require('http');
const socketIO = require('socket.io');
const fs = require('fs');
const os = require('os');

// STREAM SERVER

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

var streamRooms = {
  "name": {}
};

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

// MAIN SERVER CONNECTION

var io_client = require('socket.io-client');
var mySocket = io_client.connect("http://bid.blinkcdn.com");
mySocket.emit('connect service', "https://streamserver.blinkcdn.com", "stream");

mySocket.on('sync', function(rcvdUsers, rcvdRooms) {
  users = rcvdUsers;
  rooms = rcvdRooms;
});


/********************************/
/******* FUNCTIONALITY **********/
/********************************/

function onSignal(message, socket, destUuid, roomName) {
  room = rooms[roomName];
  var signal = JSON.parse(message);
  for (var i = 0; i < room.clients.length; i++) {
    if (room.clients[i].uuid == destUuid) {
      console.log("Sending", signal.type, " from ", socket.id, " to ", room.clients[i].socket.id)
      room.clients[i].socket.emit('signal', message);
    };
  };
  //socket.broadcast.emit('signal', message, socket.id);
}

function onDisconnect(uuid, roomName) {
    console.log(uuid, "Disconnecting");
    if(rooms[roomName]) {
        var clientsInRoom = rooms[roomName].clients
        for(var i = 0; i < clientsInRoom.length; i++) {
           console.log("i:", i);
           if (clientsInRoom[i].uuid == uuid) {
              // If this is the client, just remove them from the room
              clientsInRoom.splice(i, 1);
              rooms[roomName].clients = clientsInRoom;
           } else {
              // If this isn't the client, let them know the other client is leaving
              clientsInRoom[i].socket.emit('disconnectClient', uuid, roomName);
              console.log("Sent disconnect")
           }
        }
     }
}

function onJoin(uuid, socket, roomName) {

   if (!rooms[roomName]) {
    //If the room does not exist, create it
    console.log(socket.id, " created new room with id:", roomName);
    rooms[roomName] = {
      clients: [{"uuid": uuid, "socket": socket}]
    }
  } else if (rooms[roomName].clients.length === 1 || rooms[roomName].clients.length === 0) {
    // If rooms exist, and the most recent room only has one client,
    // add this client to the room
    clientsInThisRoom = rooms[roomName].clients
    clientsInThisRoom.push({'uuid': uuid, 'socket': socket});
    rooms[roomName].clients = clientsInThisRoom;

    clientsInThisRoom[0].socket.emit('ready', true, 2);
    clientsInThisRoom[1].socket.emit('ready', false, 2);
    console.log(socket.id, " joined the room ", roomName);

  } else if (rooms[roomName].clients.length === 2) {
    // If rooms exist, and the most recent room only has one client,
    // add this client to the room
    clientsInThisRoom = rooms[roomName].clients
    clientsInThisRoom.push({'uuid': uuid, 'socket': socket});
    rooms[roomName].clients = clientsInThisRoom;

    clientsInThisRoom[0].socket.emit('ready', true, 3);
    clientsInThisRoom[1].socket.emit('ready', false, 3);
    clientsInThisRoom[2].socket.emit('ready', false, 3);
    console.log(socket.id, " joined the room ", roomName);

  }
}

function log() {
  var array = ['Message from server:'];
  array.push.apply(array, arguments);
  socket.emit('log', array);
}
