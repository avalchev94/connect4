var tarantula = null

window.onload = function() {
  var roomName = window.location.pathname.split('/').pop()
  var playerUUID = uuidv4();

  tarantula = new Tarantula(
    new Connect4(7, 6),
    new WebSocket(`ws://localhost:8080/join?name=${roomName}&uuid=${playerUUID}`)
  )
}
