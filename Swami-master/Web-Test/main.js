function WebSocket() {
            
    if ("WebSocket" in window) {
       alert("WebSocket is supported by your Browser!");
       
       // Let us open a web socket
       var ws = new WebSocket("ws://localhost:8000/echo");
  
       ws.onopen = function() {
          
          // Web Socket is connected, send data using send()
          ws.send("Want to send a message?");
          alert("Sent message");
       };
  
       ws.onmessage = function (evt) { 
          var received_msg = evt.data;
          alert("You recieved a message");
       };
  
       ws.onclose = function() { 
          
          // websocket is closed.
          alert("You closed the connection"); 
       };
    } else {
      
       // The browser doesn't support WebSocket
       alert("WebSocket NOT supported by your Browser!");
    }
  }