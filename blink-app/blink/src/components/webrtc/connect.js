import React, { Component } from 'react';
import { connect } from 'react-redux';
import SocketIOClient from 'socket.io-client';


// Socket IO Connection
const socket = SocketIOClient.connect('https://chat.blinkcdn.com', {transports: ['websocket']});

/* user = {
  name:
  userImg:
  userID:
}*/
var user = {};


// Connects to Webrtc Signaling Server
class ConnectWebrtc extends Component {
  componentWillMount() {
    socket = io.connect();
    socket.on('created user', function(userID) {

      user.userID = userID;
        console.log("Connected");

      // Send join stream system Message
      socket.emit('join service', user.userID, 'stream', roomName);
    });

    socket.on('joined service', function(userID, serviceType, serviceAddress) {
      var engine = services[serviceType];
      engine.serviceAddress = serviceAddress;

      engine.setupService();
    });
  }
}


// Import roomName from state
const mapStateToProps = state => {
  return { roomName: state.roomName };
};

export default connect(mapStateToProps)(ConnectWebrtc);
