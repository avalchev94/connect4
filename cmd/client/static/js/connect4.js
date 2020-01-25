const Color = {
  Null: "",
  Red: "red",
  Yellow: "yellow"
}

const Message = {
  WaitingOpponent: "waiting_oponent",
  GameWon: "game_won",
  GameLost: "game_lost",
  GameDraw: "game_draw"
}

class Connect4 {
  constructor(settings, playerID) {
    this.cols = settings.cols
    this.rows = settings.rows

    this.originPlayer = playerID
    this.currentPlayer = null
    this.onMove = null
    
    this.field = []
    this.playerBox = new PlayerBox()
    this.messageBox = new MessageBox()

    // render the game UI with current game progress
    this.render(settings.gameProgress)
  }

  render(gameProgress) {
    // draw the field
    var table = document.createElement('table')
    for (var r = 0; r < this.rows; r++) {
      var row = table.insertRow()
      
      this.field[r] = []
      for (var c = 0; c < this.cols; c++) {
        var cell = row.insertCell()
        cell.onclick = this.onColumnClick.bind(this, c)

        if (gameProgress.field[r][c] != Color.Null) {
          cell.classList.add(gameProgress.field[r][c])
        }
        
        this.field[r][c] = cell
      }
    }
    document.getElementsByClassName('field')[0].appendChild(table)

    // update player and message boxes
    switch (gameProgress.state) {
    case State.Starting:
      this.messageBox.show(Message.WaitingOpponent)
      this.playerBox.hide()
      break
    case State.Running:
      this.start(gameProgress.player)
      break
    case State.Paused:
      this.messageBox.show(Message.WaitingOpponent)
      this.playerBox.hide()
      break
    case State.EndWin:
    case State.EndDraw:
      this.start(gameProgress.player)
      this.end(gameProgress.state, gameProgress.player)
      break
    }
  }

  start(player, moveRemaining) {
    this.messageBox.hide()

    this.playerBox.setPlayerNames(this.originPlayer)
    this.playerBox.show()

    this.setCurrentPlayer(player, moveRemaining)
  }

  pause() {
    this.playerBox.pause()
  }

  error(errMsg) {
    alert(`Error in game logic: ${errMsg}`)
  }

  move(player, move) {
    // change field class name to colorize the table cell
    this.field[move.row][move.col].classList.add(player)
    if (player == Color.Red) {
      this.setCurrentPlayer(Color.Yellow, 30)
    } else {
      this.setCurrentPlayer(Color.Red, 30)
    }
  }

  moveExpired(player) {
  }

  end(state, winner) {
    switch (state) {
    case State.EndWin:
      if (winner == this.originPlayer) {
        this.messageBox.show(Message.GameWon)
      } else {
        this.messageBox.show(Message.GameLost)
      }
      break
    case State.EndDraw:
      this.messageBox.show(Message.GameDraw)
      break
    }
  }

  addPlayer(player, connected) {
    this.playerBox.setPlayerConnected(player, connected)
  }

  delPlayer(player) {
  }

  setPlayerStatus(player, connected) {
    if (player == this.originPlayer) {
      throw Error("can't change connection status of origin player")
    }

    this.playerBox.setPlayerConnected(player, connected)
  }

  setCurrentPlayer(player, moveRemaining) {
    this.playerBox.setCurrentPlayer(player, moveRemaining)
    this.currentPlayer = player
  }

  playerID() {
    return this.originPlayer
  }

  onColumnClick(col) {
    if (this.currentPlayer && this.currentPlayer == this.originPlayer) {
      var move = {col: col, row: -1}
      
      for (var row = this.field.length - 1; row >= 0 ; row--) {
        var classes = this.field[row][move.col].classList
        if (!classes.contains('red') && !classes.contains('yellow')) {
          move.row = row
          break;
        }
      }

      if (move.row == -1) {
        alert("column is full, choose another")
        return
      }

      this.onMove(move)
    }
  }
}

class MessageBox {
  constructor() {
    this.box = document.getElementsByClassName('content message')[0]
    this.messages = new Map([
      [Message.WaitingOpponent, document.getElementsByClassName('waiting-opponent')[0]],
      [Message.GameWon, document.getElementsByClassName('end win')[0]],
      [Message.GameLost, document.getElementsByClassName('end loss')[0]],
      [Message.GameDraw, document.getElementsByClassName('end draw')[0]],
    ])
  }

  hide() {
    this.box.setAttribute("hidden", "")
    this.messages.forEach(function(value){
      value.setAttribute("hidden", "")
    })
  }

  show(msg) {
    this.hide()

    this.box.removeAttribute("hidden")
    this.messages.get(msg).removeAttribute("hidden")
  }
}

class PlayerBox {
  constructor() {
    this.box = document.getElementsByClassName('game-info')[0]
    this.pauseIcon = document.getElementsByClassName('pause-icon')[0]
    this.players = new Map([
      [Color.Red, new Player(Color.Red)],
      [Color.Yellow, new Player(Color.Yellow)],
    ])
    this.timer = new Timer(0)
  }

  hide() {
    this.box.style.display = "none"
  }

  show() {
    this.box.style.display = "flex"
  }

  togglePauseIcon(state) {
    this.pauseIcon.style.visibility = state ? "visible" : "hidden"
  }

  pause() {
    this.togglePauseIcon(true)
    this.timer.stop()
  }

  setPlayerConnected(player, connected) {
    this.players.get(player).setConnected(connected)
  }

  setPlayerNames(player) {
    this.players.forEach(function(value, key){
      if (player == key) {
        value.setName("You")
      } else {
        value.setName("Opponent")
      }
    })
  }

  setCurrentPlayer(player, moveRemaining) {
    this.players.get(Color.Red).setCurrent(player == Color.Red)
    this.players.get(Color.Yellow).setCurrent(player == Color.Yellow)
    
    this.togglePauseIcon(false)
    this.timer.reset(moveRemaining)
  }
}

class Player {
  constructor(color) {
    this.color = color
    this.player = document.getElementsByClassName(`player ${color}`)[0]
    this.name = document.getElementById(`${color}-name`)
  }

  setName(name) {
    this.name.innerHTML = name
  }

  setConnected(state) {
    console.log(this.color, ' connection state:', state)
  }

  setCurrent(state) {
    if (state) {
      this.player.classList.add('on-move')
    } else {
      this.player.classList.remove('on-move')
    }
  }
}

class Timer {
  constructor(duration) {
    this.timerBox = document.getElementsByClassName('timer')[0]
    this.timerText = this.timerBox.getElementsByClassName('remaining')[0]
    this.running = false
    this.duration = duration
    this.remaining = duration
    this.intervalID = -1

    this.setDuration(duration)
  }

  setDuration(duration) {
    if (!this.running) {
      this.duration = duration
      this.remaining = duration
      this.timerText.innerHTML = duration
    }
  }

  start() {
    if (this.running || this.remaining <= 0) {
      return
    }

    this.running = true
    this.intervalID = setInterval(()=>{
      this.timerText.innerHTML = --this.remaining
      if (this.remaining == 0) {
        this.stop()
      }
    }, 1000)
  }

  stop() {
    if (this.running || this.intervalID != -1) {
      clearInterval(this.intervalID)
      this.running = false
      this.intervalID = -1
    }
  }

  reset(duration) {
    if (this.running) {
      this.stop()
    }

    this.setDuration(duration)
    this.start()
  }
}