
const getQueueHandler = async (req, res, { Queue }) => {
    try{
        const queue = await Queue.find({});
        res.status(200).json(queue)
    }catch(e) {

    }
}

const postQueueHandler = async (req, res, { Queue, con }) => {
    try {
        const temp = JSON.parse(req.headers["x-user"]);
        var userID = temp.id;
        var { topicID, quizAnswer } = req.body
        if (!topicID) {
            res.status(400).send("Must provide a topic ID")
            return;
        }

        if (quizAnswer === undefined) {
            res.status(400).send("Must provide a quiz answer")
            return;
        }
        const queuedAt = new Date();
        const queueUsers = await Queue.find({topicID: topicID, quizAnswer: !quizAnswer});
        var dt = new Date().getTime();
        var roomId = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = (dt + Math.random()*16)%16 | 0;
            dt = Math.floor(dt/16);
            return (c=='x' ? r :(r&0x3|0x8)).toString(16);
        });
        let queue = {
            userID,
            topicID,
            quizAnswer,
            roomId,
            queuedAt
        };
        if (queueUsers.length == 0) {
            const query = new Queue(queue)
            query.save((err, newQueue) => {
                if (err) {
                    res.status(500).send(err.toString())
                    return;
                }
                res.setHeader("content-type", "text/plain");
                res.status(201).send(roomId)
            })
        } else {
            let match = queueUsers[0];
            await Queue.remove({userID: match.userID})
            let roomId = match.roomId
            res.setHeader("content-type", "text/plain")
            res.status(201).send(roomId)

        }
    } catch(e) {
        res.status(500).send(e.toString());
    }
}

const deleteQueueHandler = async (req, res, { Queue }) => {
    try {
        const temp = JSON.parse(req.headers["x-user"]);
        var userID = temp.id;  
        const queueUsers = await Queue.find({userID: userID}); 
        if (queueUsers.length == 0) {
            res.status(500).send("You aren't in a queue right now")
            return;
        } 
        await Queue.remove({userID: userID})
        res.setHeader("content-type", "text/plain");
        res.status(201).send("Removed from queue")

    } catch(e) {
        res.status(500).send(e.toString());
    }
}

module.exports = { postQueueHandler, deleteQueueHandler, getQueueHandler }