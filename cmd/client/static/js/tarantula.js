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
  GamePaused: "game_paused",
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
//    pause() // game pauses
//    move(player, move) // player has made move
//    moveExpired(player, move) // player move time has expired
//    end(state, player) // game ended with state
//    addPlayer(player) // new player added, connection status false by default
//    delPlayer(player) // player left
//    setPlayerStatus(player, connected) // set player current connection status
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
      this.game.start(msg.payload.starting, msg.payload.moveRemaining)
      break
    case MessageType.GameEnded:
      this.game.end(msg.payload.state, msg.payload.winner)
      break
    case MessageType.GamePaused:
      this.game.pause()
    case MessageType.PlayerMove:
      this.game.move(msg.payload.player, msg.payload.move)
      break
    case MessageType.PlayerMoveExpired:
      this.game.moveExpired(msg.payload.player)
    case MessageType.PlayerJoined:
      this.game.addPlayer(msg.payload.player)
      break
    case MessageType.PlayerLeft:
      this.game.delPlayer(msg.payload.player)
      break
    case MessageType.PlayerConnected:
      this.game.setPlayerStatus(msg.payload.player, true)
    case MessageType.PlayerDisconnected:
      this.game.setPlayerStatus(msg.payload.player, false)
    }
  }
}