const HTTPS_PORT = 443;
const HTTP_PORT = 80;

const nodeStatic = require('node-static');
const https = require('https');
const http = require('http');
const socketIO = require('socket.io');
const fs = require('fs');
const os = require('os');

////////////
var users = {
  // uuid: {}
};

var sockets = {};
var services = {
  bid: {
    address: null,
    socket: null,
  },
  stream: {
    address: null,
    socket: null
  }
};

var rooms = {
  // roomName: {
  //   name: "name", // Identifier of the room
  //   services: [], // Array of allowed services
  //   users: {} // Dictionary of members allowed (for easy pulling)
  // },
}


/************ SERVER SETUP *************/
const certOptions = {
  key: fs.readFileSync('certs/key.pem'),
  cert: fs.readFileSync('certs/cert.pem')
}

var fileServer = new(nodeStatic.Server)();
var app = https.createServer(certOptions, function(req, res) {
  fileServer.serve(req, res);
}).listen(HTTPS_PORT);

var io = socketIO.listen(app);

io.sockets.on('connection', function(socket) {

  socket.on('create user', function(user, roomName) {
    createUser(user, roomName, socket);
  })

  socket.on('join service', function(userID, serviceType, roomName) {
    setupService(userID, serviceType, roomName, socket);
    console.log("Joined System.", roomName);
  });

  socket.on('place bid', function(bidValue) {
    console.log(bidValue);
  });

})

/************ SERVICES SERVER SETUP *************/
var serviceFileServer = new(nodeStatic.Server)();
var serviceApp = http.createServer(function(req, res) {
  serviceFileServer.serve(req, res);
}).listen(HTTP_PORT);

var serviceIo = socketIO.listen(serviceApp);

serviceIo.sockets.on('connection', function(socket) {

  socket.on('connect service', function(serviceAddress, serviceType) {
    services[serviceType].address = serviceAddress;
    services[serviceType].socket = socket;
    console.log("Connected Service:", serviceType, serviceAddress);
    syncUpdateService(serviceType);
  });

  socket.on('sync user request', function() {
    socket.emit('sync', users, rooms);
  });

  socket.on('sync', function(rcvdUsers, rcvdRooms) {
    users = rcvdUsers;
    rooms = rcvdRooms;
    updateAllServices();
  });
});

console.log("Connected.");
/******** FUNCTIONS *********/

function createUser(user, roomName, socket) {
  var newUser = {
    userID: uuid(),
    name: user.name,
    userImg: user.userImg
  }

  // Add user to the array of users
  sockets[newUser.userID] = socket;
  users[newUser.userID] = newUser;

  // Add user to room
  if(!rooms[roomName]) {
    rooms[roomName] = {
      users: {},
    }

    rooms[roomName].roomName = roomName;
    rooms[roomName].users[newUser.userID] = newUser
  } else {
    rooms[roomName].users[newUser.userID] = newUser;
  }

  socket.emit('created user', newUser.userID, newUser.name);
}

// Create user for system
function setupService(userID, serviceType, roomName, socket) {

  // Set the proper server address
  var serviceAddress;

  // Add service structure
  if (!rooms[roomName].hasOwnProperty('services')) {
    rooms[roomName].services = {};
  }

  if (!rooms[roomName].services.hasOwnProperty(serviceType)) {
    rooms[roomName].services[serviceType] = createService(serviceType);
  }

  if (services[serviceType]) {
    serviceAddress = services[serviceType].address;
    socket.emit('joined service', userID, serviceType, serviceAddress);
    console.log("Joined service:", userID, serviceType, serviceAddress);
  } else {
    console.log("Service to setup not found.");
  }

  syncUpdateService(serviceType);
}

function syncUpdateService(serviceType) {
  if(services[serviceType].socket) {
    services[serviceType].socket.emit('sync', users, rooms);
  } else {
    console.log("Failed to update service. Please check service type:", serviceType);
  }
}

function updateAllServices() {
  for (service in services) {
    syncUpdateService(service);
  }
}

/****** HELPER FUNCTIONS ******/

function createService(serviceType) {

  if (serviceType == "bid") {
    return {
      lots: [{
        itemName: "Mona Lisa",
        description: "Does it need a description?",
        itemImg: "items/monalisa.jpg"
        },
        {
        itemName: "Rolex Daytona 116500",
        description: "High end fancy watch.",
        itemImg: "items/rolex.jpg"
        }],
      currentLotNumber: 0,
      highestBid: 0,
      bidCount: 0,
      bids: [],
      hasStartedBid: false,
    }
  }

}

function uuid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
  }

  return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
}
