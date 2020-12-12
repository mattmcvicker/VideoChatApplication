function startCall() {
    const socket = io('/')
    const videoGrid = document.getElementById('video-grid')

    // create a new Peer using the Peerjs library
    // undefined is where the userID goes, if undefined,
    // then the Peer autogenerates it
    const myPeer = new Peer(undefined, {
        host: '/',
        port: '3001'
    })

    // creates a video html object and mutes it (only mutes it for 
    // ourselves, so prevent echo)
    const myVideo = document.createElement('video')
    myVideo.muted = true

    // keeps track of all the peers in the call
    const peers = {}

    // gets the video and audio from the device currently using this
    navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true
    }).then(stream => {
        // we wrote this function to connect the stream to the correct
        // things on the HTML/EJS side
        addVideoStream(myVideo, stream)

        // this makes it so that the myPeer object actually listens for
        // for calls from other myPeer objects
        myPeer.on('call', call => {
            console.log("recieved a call")
            call.answer(stream)

            // this makes it so that the call reciever also gets the video
            // stream from the call sender
            const video = document.createElement('video')
            call.on('stream', userVideoStream => {
                addVideoStream(video, userVideoStream)
            })
        })

        // when a new user is connected, connect the new user's stream
        // given their userId
        socket.on('user-connected', userId => {
            console.log(userId, "has joined")
            connectToNewUser(userId, stream)
        })
    })

    socket.on('user-disconnected', userId => {
        if (peers[userId]) peers[userId].close
    })

    // when myPeer is created, sends out a 'join-room' event with the
    // given id
    myPeer.on('open', id => {
        socket.emit('join-room', ROOM_ID, id)
    })

    // connects to the new user using Peer's call method
    function connectToNewUser(userId, stream) {
        // uses the myPeer object to call the other user
        console.log("calling", userId)
        const call = myPeer.call(userId, stream)
        const video = document.createElement('video')

        // when we get the stream back from the user we are calling,
        // adds that video stream to the video grid
        call.on('stream', userVideoStream => {
            console.log("We recieved a stream back")
            addVideoStream(video, userVideoStream)
        })

        // when a user closes the call, their video is removed
        call.on('close', () => {
            video.remove()
        })

        // puts the call in the peers object mapped to the correct
        // userId
        peers[userId] = call
    }

    // takes in a video html object, and connects a given stream to it
    function addVideoStream(video, stream) {
        video.srcObject = stream
        // when the stream is loaded into the video object, start
        //playing the video
        video.addEventListener('loadedmetadata', () => {
            video.play()
        })
        videoGrid.append(video)
    }
}