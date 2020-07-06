function sendData() {
    var about = document.getElementById("about").value;
    var rollno = document.getElementById("rollno").value;
    if (about.trim().length > 1){
      var socketMsg= {"rollno" : rollno.toString(), "about" : about};
      socketSend(socketMsg);
    }
  }
    var sockEcho = null;
    var oldRollNo = null;
    var wsServerAddress = "ws://127.0.0.1:4050";
  
    window.onload = function() {
      // Connect the WebSocket to the server and register callbacks on it.
      sockEcho = new WebSocket(wsServerAddress + "/student");
      sockEcho.onopen = function() {
        console.log("connected");
      }
  
      sockEcho.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
      }
    };
  
    function socketSend(msg) {
      if (sockEcho != null && sockEcho.readyState == WebSocket.OPEN) {
        sockEcho.send(JSON.stringify(msg));
      } else {
        console.log("Socket isn't OPEN");
      }
    }
  