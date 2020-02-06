var tarantula = null
var hostname = document.location.hostname

window.onload = function() {
  var name = document.URL.split('/').pop()
  fetch(`http://${hostname}:8080/rooms/${name}/settings`, {
    method: 'GET',
    credentials: 'include'
  })
  .then((resp) => {
    if (resp.ok) {
      return resp.json()
    } else {
      throw resp.text()
    }
  })
  .then((resp) => {
    tarantula = new Tarantula(
      new Connect4(resp.settings, resp.player),
      new this.WebSocket(`ws://${hostname}:8080/rooms/${name}/connect`)
    )
  })
  .catch((error) => {
    this.alert(error)
  })


}