const STATE = {
  Starting: 0,
  Running: 1,
  EndDraw: 2,
  EndWin: 3
}

class Tarantula {
  constructor(game, socket) {
    this.game = game
    this.game.onaction = this.onMove.bind(this)

    this.socket = socket
    this.socket.onmessage = this.onMessage.bind(this)

    this.game.render()
  }

  onMove(move) {
    var msg = {
      Move: move,
      Player: this.game.playerID()
    }

    this.game.move(msg.Player, msg.Move)
    this.socket.send(JSON.stringify(msg))
  }

  onMessage(event) {
    var msg = JSON.parse(event.data)
    switch (msg.State) {
    case STATE.Starting:
      this.game.start(msg.Player)
      break
    case STATE.Running:
      this.game.move(msg.Player, msg.Move)
      break
    case STATE.EndDraw:
    case STATE.EndWin:
      this.game.end(msg.State, msg.Player)
      break
    }
  }
}