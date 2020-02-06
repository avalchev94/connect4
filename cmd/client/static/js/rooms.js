var hostname = document.location.hostname

function updateRooms() {
  var table = document.getElementsByClassName('rooms-table')[0]
  var noRooms = document.getElementsByClassName('no-rooms')[0]

  // delete rooms
  while (table.rows.length > 1) {
    table.deleteRow(-1)
  }

  // fetch new rooms
  fetch(`http://${hostname}:8080/rooms`).then(function(resp){
    resp.json().then(function(rooms){
      rooms.forEach(function(room){
        var row = table.insertRow()
        row.insertCell().innerText = room.name
        row.insertCell().innerText = room.players + '/2'
        row.insertCell().innerText = room.game
        var joinButton = row.insertCell()
        joinButton.className = 'join_button'
        joinButton.onclick = joinRoom.bind(null, room)
      })
    }).then(function(){
      noRooms.hidden = table.rows.length > 1
    })
  })

}

function createRoom() {
  var body = {
    name: document.getElementsByName('room_name')[0].value,
    game: document.getElementsByName('game_name')[0].value
    // settigs
  }
  
  fetch(`http://${hostname}:8080/rooms/new`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
  .then((resp) => {
    if (resp.status == 201) {
      joinRoom(body)
    } else {
      throw resp.text()
    }
  })
  .catch((error) => {
    error.then((message) => {
      alert(message)
    })
  })

}

function joinRoom(room) {
  fetch(`http://${hostname}:8080/rooms/${room.name}/join`, {
    method: 'POST',
    credentials: 'include'
  })
  .then((resp) => {
    if (resp.ok) {
      document.location.href = `/${room.game}/${room.name}`
    } else {
      throw resp.text()
    }
  })
  .catch((error) => {
    error.then((message) => {
      alert(message)
    })
  })
}

window.onload = function() {
  updateRooms()
  //setInterval(updateRooms, 5000)
}


