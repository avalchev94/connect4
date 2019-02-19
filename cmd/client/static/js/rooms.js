function updateRooms() {
  var table = document.getElementsByClassName('rooms-table')[0]
  var noRooms = document.getElementsByClassName('no-rooms')[0]

  // delete rooms
  while (table.rows.length > 1) {
    table.deleteRow(-1)
  }

  // fetch new rooms
  fetch('http://localhost:8080/rooms').then(function(resp){
    resp.json().then(function(rooms){
      rooms.forEach(function(room){
        var row = table.insertRow()
        row.insertCell().innerText = room.Name
        row.insertCell().innerText = room.Players + '/2'
        row.insertCell().innerText = room.Game
        var joinButton = row.insertCell()
        joinButton.className = 'join_button'
        joinButton.onclick = joinRoom.bind(null, room.Name)
      })
    }).then(function(){
      noRooms.hidden = table.rows.length > 1
    })
  })

}

function createRoom() {
  var roomName = document.getElementsByName('room_name')[0].value
  
  fetch('http://localhost:8080/new?name='+roomName).then(function(resp){
    if (resp.status == 201) {
      window.location.href = "/connect4/"+roomName
    }
  })  
}

function joinRoom(roomName) {
  window.location.href = "/connect4/"+roomName
} 

window.onload = function() {
  updateRooms()
  //setInterval(updateRooms, 5000)
}


