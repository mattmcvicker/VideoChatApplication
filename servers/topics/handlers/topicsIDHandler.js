const getSpecificTopicsHandler = async (req, res, { Topic }) => {
    try {
        const id = req.params.topicID;
        const topic = await Topic.find({_id: id}); 
        if (topic.length < 1) {
            res.status(500).send("Topic does not exist")
        } 
        res.setHeader("content-type", "application/json");
        res.json(topic);  
    } catch(e) {
        res.status(500).send(e.toString())   
    }
}

const patchSpecificTopicsHandler = async (req, res, { Topic }) => {
    try {
        const temp = req.headers["x-user"];
        var userID = JSON.parse(temp).id

        var topicID = req.body.topicID;
        const topic = await Topic.find({_id: topicID})[0]; 
        if (topic.votes.contains(userID)) {
            topic.votes.delete(userID) //delete user if already voted
        } else {
            topic.votes.add(userID)
        }

        const query = topic;
        query.save((err, updatedTopic) => {
            if (err) {
                res.status(500).send(err.toString())
                return;
            }
            res.setHeader("content-type", "application/json");
            res.status(201).json(updatedTopic)
    })  
    } catch(e) {
        res.status(500).send(e.toString())   
    }
}

const deleteSpecificTopicsHandler = async (req, res, { Topic }) => {
    try {
        const temp = req.headers["x-user"];
        var userID = JSON.parse(temp).id


        const topicID = req.params.topicID;
        const topic = await Topic.find({_id: topicID})[0]; 
        if (topic.length < 1) {
            res.status(500).send("Topic does not exist")
        } 
        if (userID != topic.creator) {
            res.status(500).send("You must be the creator to delete a topic")
        }
        await Topic.remove({_id: topicID});
        res.setHeader("content-type", "text/plain");
        res.send("Topic deleted") 
    } catch(e) {
        res.status(500).send(e.toString())   
    }
}


module.exports = { getSpecificTopicsHandler, patchSpecificTopicsHandler, deleteSpecificTopicsHandler }