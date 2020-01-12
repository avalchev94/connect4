const State = {
  Starting: "starting",
  Running: "running",
  Paused: "paused",
  EndDraw: "end_draw",
  EndWin: "end_win"
}

const MessageType = {
	GameStarting: "game_starting",
	GameEnded: "game_ended",
  PlayerMove: "player_move",
  PlayerMoveExpired: "player_move_expired",
	PlayerJoined: "player_joined",
  PlayerLeft: "player_left",
  PlayerConnected: "player_connected",
  PlayerDisconnected: "player_disconnected"
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