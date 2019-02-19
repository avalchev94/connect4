const COLOR = {
  Null: 0,
  Red: 1,
  Yellow: 2
}

class Connect4 {
  constructor(cols, rows) {
    this.cols = cols
    this.rows = rows
    this.originPlayer = null
    this.currentPlayer = null
    this.onaction = null
    
    this.field = []
    this.players = new Map([
      [COLOR.Red, document.getElementsByClassName('player red')[0]],
      [COLOR.Yellow, document.getElementsByClassName('player yellow')[0]]
    ])
  }

  render() {
    var table = document.createElement('table')
    for (var r = 0; r < this.rows; r++) {
      var row = table.insertRow()
      
      this.field[r] = []
      for (var c = 0; c < this.cols; c++) {
        var cell = row.insertCell()
        cell.onclick = this.onColumnClick.bind(this, c)
        
        this.field[r][c] = cell
      }
    }

    document.getElementsByClassName('field')[0].appendChild(table)
  }

  start(player) {
    this.originPlayer = player

    // red always start first
    this.setCurrentPlayer(COLOR.Red)
  }

  move(player, move) {
    debugger
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

  end(state, player) {
    switch (state) {
    case STATE.EndWin:
      alert(player + 'wins.')
      break
    case STATE.EndDraw:
      alert('draw')  
      break
    }
  }

  setCurrentPlayer(player) {
    if (player == COLOR.Red) {
      this.players.get(COLOR.Red).classList.add('on-move')
      this.players.get(COLOR.Yellow).classList.remove('on-move')
    } else {
      this.players.get(COLOR.Red).classList.remove('on-move')
      this.players.get(COLOR.Yellow).classList.add('on-move')   
    }

    this.currentPlayer = player
  }

  playerID() {
    return this.originPlayer
  }

  onColumnClick(col) {
    if (this.currentPlayer == this.originPlayer) {
      this.onaction({col: col, row: -1})
    }
  }
}

var tarantula = null

window.onload = function() {
  var roomName = window.location.pathname.split('/').pop();

  var game = new Connect4(7, 6)
  var socket = new WebSocket('ws://localhost:8080/join?name='+roomName)

  tarantula = new Tarantula(game, socket)
}