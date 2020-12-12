const express = require('express')
const app = express()
const server = require('http').Server(app)
const io = require('socket.io')(server)
const { v4: uuidV4 } = require('uuid') //creates a function that creates ids

app.set('view engine', 'ejs')
app.use(express.static('public'))

// the default path, redirects you do randomly generated room
app.get('/', (req, res) => {
    res.redirect(`/${uuidV4()}`)
})

// renders a room with the id in the path
app.get('/:room', (req, res) => {
    res.render('room', { roomId: req.params.room })
})

// io searches for connection, and executes once it is connection
io.on('connection', socket => {
    // websocket waits for a join-room event, then has the user join the room
    // with the given room ID
    socket.on('join-room', (roomId, userId) => {
        socket.join(roomId)
        socket.to(roomId).broadcast.emit('user-connected', userId)

        socket.on('disconnect', () => {
            socket.to(roomId).broadcast.emit('user-disconnected', userId)
        })
    })
})

server.listen(3000)