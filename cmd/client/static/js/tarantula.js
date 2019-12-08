const State = {
  Starting: 0,
  Running: 1,
  Paused: 2,
  EndDraw: 3,
  EndWin: 4
}

const MessageType = {
	GameStarting: 0,
	GameEnded: 1,
	PlayerMove: 2,
	PlayerJoined: 3,
	PlayerLeft: 4
}

// interface Game {
//    render() // render the UI
//    start(player) // game starts, first is given player
//    move(player, move) // player has made move
//    end(state, player) // game ended with state
// }

class Tarantula {
  constructor(game, socket) {
    this.game = game
    this.game.onaction = this.onMove.bind(this)

    this.socket = socket
    this.socket.onmessage = this.onMessage.bind(this)
    this.socket.onerror = function(event) {
      console.log(event)
    };

    this.game.render()
  }

  onMove(move) {
    var msg = {
      type: MessageType.PlayerMove,
      payload: {
        player: this.game.playerID(),
        move: move
      }
    }

    this.game.move(this.game.playerID(), move)
    this.socket.send(JSON.stringify(msg))
  }

  onMessage(event) {
    var msg = JSON.parse(event.data)

    switch (msg.type) {
    case MessageType.GameStarting:
      this.game.start(msg.payload.starting)
      break
    case MessageType.GameEnded:
      this.game.end(msg.payload.state, msg.payload.winner)
      break
    case MessageType.PlayerMove:
      this.game.move(msg.payload.player, msg.payload.move)
      break
    }
  }
}