var wsServerAddress = "ws://127.0.0.1:4050";
var socketConn = null;
var originalData = {}
window.onload = function() {
  var table = document.getElementById("outputTable");
  var rows = table.getElementsByTagName('tr');
  socketConn = new WebSocket(wsServerAddress + "/teacher");
  socketConn.onmessage = function(e) {
    newData = JSON.parse(e.data)
    for (var key in newData) {
      if(originalData[key] != undefined){
        rowIndex = table.rows[originalData[key]].cells;
        rowIndex[1].innerHTML = newData[key]['wordCount']
        rowIndex[2].innerHTML = newData[key]['charCount']
        rowIndex[3].innerHTML = newData[key]['actualMessage']
      }else{
        appendRowsToTable(table,newData,key)
      }
    }
  }

  function appendRowsToTable(tableRef,data,key){
    var row = tableRef.insertRow(-1)
    var rollNo = row.insertCell(0);
    var totalWords = row.insertCell(1);
    var totalCharacters = row.insertCell(2);
    var actualMessage = row.insertCell(3);
    rollNo.innerHTML = data[key]['rollNo']
    totalWords.innerHTML = data[key]['wordCount']
    totalCharacters.innerHTML = data[key]['charCount']
    actualMessage.innerHTML = data[key]['actualMessage']
    originalData[key] = tableRef.rows.length - 1;
  }


  socketConn.onopen = function() {
    console.log("connected");
  }

  socketConn.onclose = function(e) {
    console.log("connection closed (" + e.code + ")");
  }
};