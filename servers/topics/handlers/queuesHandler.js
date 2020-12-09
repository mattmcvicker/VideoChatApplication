
const postQueueHandler = async (req, res, { Queue, con }) => {
    try {
        const temp = JSON.parse(req.headers["x-user"]);
        var userID = temp.id;
        var { quizID, quizAnswer } = req.body
        if (!quizID) {
            res.status(400).send("Must provide a quiz ID")
            return;
        }

        if (!quizAnswer) {
            res.status(400).send("Must provide a quiz answer")
            return;
        }
        const queuedAt = new Date();
        const queueUsers = await Queue.find({quizID: quizID, quizAnswer: !quizAnswer});
        let queue = {
            userID,
            quizID,
            quizAnswer,
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
                res.status(201).send("User put in queue successfully")
            })
        } else {
            let match = queueUsers[0];
            await Queue.remove({userID: match.userID})
            let connectionInfo = {
                user_1ID: userID,
                user_2ID: match.userID
            }
            res.setHeader("content-type", "application/json")
            res.status(201).json(connectionInfo)

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

module.exports = { postQueueHandler, deleteQueueHandler }