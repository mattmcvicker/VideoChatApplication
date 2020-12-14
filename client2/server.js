const express = require('express')
const app = express()
const https = require('https')
const fs = require('fs')
// const server require('http').Server(app)
// const server = https.Server(app)
var privateKey = fs.readFileSync('/etc/letsencrypt/live/kenmasumoto.me/privkey.pem')
var certificate = fs.readFileSync('/etc/letsencrypt/live/kenmasumoto.me/fullchain.pem')
var credentials = {key: privateKey, cert: certificate}
var httpsServer = https.createServer(credentials, app)
const ExpressPeerServer = require('peer').ExpressPeerServer
const peerServer = ExpressPeerServer(httpsServer, { 
    port: 443,
    ssl: {
        key: fs.readFileSync('/etc/letsencrypt/live/kenmasumoto.me/privkey.pem'),
        cert: fs.readFileSync('/etc/letsencrypt/live/kenmasumoto.me/fullchain.pem')
    }
})
app.use('/peerjs', peerServer)
const io = require('socket.io')(httpsServer)
const { v4: uuidV4 } = require('uuid') //creates a function that creates ids

app.set('view engine', 'ejs')
app.use(express.static('public'))

// the default path, redirects you do randomly generated room
app.get('/', (req, res) => {
    //var roomId = req.headers('roomId')
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


httpsServer.listen(443)
//httpsServer.listen(9000)

//server.listen(443)
