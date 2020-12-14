const getTopicsHandler = async (req, res, { Topic }) => {
    try {
        const allTopics = await Topic.find({});  
        res.setHeader("content-type", "application/json");
        res.json(allTopics);  
    } catch(e) {
        res.status(500).send(e.toString())   
    }
}

const postTopicsHandler = async (req, res, { Topic }) => {
    try {
        const temp = req.headers["x-user"];
        var ID = temp;
        ID = ID.split("id: ")[1].split("}")[0];
        var { name } = req.body
        if (!name) {
            res.status(400).send("Must provide a name for the Topic")
            return;
        }
        
        const createdAt = new Date();
        const editedAt = new Date();
        var initialVotes = 0;
        var topic = {};
        // var userEmail = "";

        topic = {
            name,
            votes: [],
            createdAt,
            creator: ID,
            editedAt
        }
        // res.status(500).send(channel)
        const query = new Topic(topic);
        query.save((err, newTopic) => {
            if (err) {
                res.status(500).send(err.toString())
                return;
            }
            res.setHeader("content-type", "application/json");
            res.status(201).json(newTopic)
    })  
    } catch(e) {
        res.status(500).send(e.toString())   
    }
}

module.exports = { getTopicsHandler, postTopicsHandler }