const COLOR = {
  Null: 0,
  Red: 1,
  Yellow: 2
}

const MESSAGE = {
  WaitingOpponent: 0,
  GameWon: 1,
  GameLost: 2,
  GameDraw: 3,
  ConnectionLost: 4
}

class Connect4 {
  constructor(settings, playerID) {
    this.cols = settings.cols
    this.rows = settings.rows

    this.originPlayer = playerID
    this.currentPlayer = null
    this.onaction = null
    
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
        
        if (gameProgress.field[r][c] == COLOR.Red) {
          cell.classList.add('red')
        } else if (gameProgress.field[r][c] == COLOR.Yellow) {
          cell.classList.add("yellow")
        }
        
        this.field[r][c] = cell
      }
    }
    document.getElementsByClassName('field')[0].appendChild(table)

    // update player and message boxes
    switch (gameProgress.state) {
    case State.Starting:
      this.messageBox.show(MESSAGE.WaitingOpponent)
      this.playerBox.hide()
      break
    case State.Running:
      this.start(gameProgress.player)
      break
    case State.Paused:
      this.messageBox.show(MESSAGE.WaitingOpponent)
      this.playerBox.hide()
      break
    case State.EndWin:
    case State.EndDraw:
      this.start(gameProgress.player)
      this.end(gameProgress.state, gameProgress.player)
      break
    }
  }

  start(player) {
    this.setCurrentPlayer(player)
    this.playerBox.setPlayerNames(this.originPlayer)

    // show playerBox and hide messages
    this.playerBox.show()
    this.messageBox.hide()
  }

  move(player, move) {
    if (move.row == -1) {
      for (var row = this.field.length - 1; row >= 0 ; row--) {
        var classes = this.field[row][move.col].classList
        if (!classes.contains('red') && !classes.contains('yellow')) {
          move.row = row
          break;
        }
      }
    }

    // change field class name to colorize the table cell
    var cell = this.field[move.row][move.col].classList
    if (player == COLOR.Red) {
      cell.add('red')
      this.setCurrentPlayer(COLOR.Yellow)
    } else {
      cell.add('yellow')
      this.setCurrentPlayer(COLOR.Red)
    }
  }

  end(state, winner) {
    switch (state) {
    case State.EndWin:
      if (winner == this.originPlayer) {
        this.messageBox.show(MESSAGE.GameWon)
      } else {
        this.messageBox.show(MESSAGE.GameLost)
      }
      break
    case State.EndDraw:
      this.messageBox.show(MESSAGE.GameDraw)
      break
    }
  }

  setCurrentPlayer(player) {
    this.playerBox.setCurrentPlayer(player)
    this.currentPlayer = player
  }

  playerID() {
    return this.originPlayer
  }

  onColumnClick(col) {
    if ( this.currentPlayer 
      && this.originPlayer
      && this.currentPlayer == this.originPlayer
    ) {
      this.onaction({col: col, row: -1})
    }
  }
}

class MessageBox {
  constructor() {
    this.messageBox = document.getElementsByClassName('content message')[0]
    this.messages = new Map([
      [MESSAGE.WaitingOpponent, document.getElementsByClassName('waiting-opponent')[0]],
      [MESSAGE.GameWon, document.getElementsByClassName('end win')[0]],
      [MESSAGE.GameLost, document.getElementsByClassName('end loss')[0]],
      [MESSAGE.GameDraw, document.getElementsByClassName('end draw')[0]],
      [MESSAGE.ConnectionLost, document.getElementsByClassName('disconnected')[0]]
    ])
  }

  hide() {
    this.messageBox.setAttribute("hidden", "")
    this.messages.forEach(function(value){
      value.setAttribute("hidden", "")
    })
  }

  show(msg) {
    this.hide()

    this.messageBox.removeAttribute("hidden")
    this.messages.get(msg).removeAttribute("hidden")
  }
}

class PlayerBox {
  constructor() {
    this.playersBox = document.getElementsByClassName('game-info')[0]
    this.players = new Map([
      [COLOR.Red, document.getElementsByClassName('player red')[0]],
      [COLOR.Yellow, document.getElementsByClassName('player yellow')[0]]
    ])
    this.playersNames = new Map([
      [COLOR.Red, document.getElementById('red-name')],
      [COLOR.Yellow, document.getElementById('yellow-name')]
    ])
    this.timer = new Timer(30)
  }

  hide() {
    this.playersBox.style.display = "none"
  }

  show() {
    this.playersBox.style.display = "flex"
  }

  setPlayerNames(player) {
    this.playersNames.forEach(function(value, key){
      value.innerHTML = key == player ? "You" : "Opponent"
    })
  }

  setCurrentPlayer(player) {
    if (player == COLOR.Red) {
      this.players.get(COLOR.Red).classList.add('on-move')
      this.players.get(COLOR.Yellow).classList.remove('on-move')
    } else {
      this.players.get(COLOR.Red).classList.remove('on-move')
      this.players.get(COLOR.Yellow).classList.add('on-move')
    }
    
    this.timer.restart()
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

  restart() {
    if (this.running) {
      this.stop()
    }

    this.setDuration(this.duration)
    this.start()
  }
}