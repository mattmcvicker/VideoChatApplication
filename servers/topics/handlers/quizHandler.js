// define quiz handlers

// returns a specific quiz
const getSpecificQuizHandler = async (req, res, { Topic, Quiz }) => {
    try {
        // check that the user is logged in
        const temp = req.headers["x-user"];
        var userID = temp;
        userID = userID.split("id: ")[1].split("}")[0];
        if (!userID) {
            res.status(403).send("Not authorized: user must be logged in")
        }

        const topicid = req.params.topicID;
        const topic = await Topic.find({_id: topicid}); 
        if (topic.length < 1) {
            res.status(404).send("Not found: Topic does not exist");
            return;
        }
    
        // get quiz for the topic
        const quiz = await Quiz.find({topicid: topicid});
        if (quiz.length < 1) {
            res.status(404).send("Not found: Quiz deos not exist");
            return;
        }
        res.status(200).setHeader("content-type", "application/json").json(quiz);
    } catch (e) {
        res.status(500).send("Internal server error: " + e.toString());
    }
}

// creates a new quiz for a topic
const postQuizHandler = async (req, res, { Topic, Quiz }) => {
    try {
        // check that the user is logged in
        const temp = req.headers["x-user"];
        var userID = temp;
        userID = userID.split("id: ")[1].split("}")[0];
        if (!userID) {
            res.status(403).send("Not authorized: user must be logged in")
        }

        const topicid = req.params.topicID;
        const topic = await Topic.find({_id: topicid}); 
        if (topic.length < 1) {
            res.status(404).send("Not found: Topic does not exist");
            return;
        }

        // check if quiz already exists
        const quiz = await Quiz.find({topicid: topicid});
        if (quiz.length >= 1) {
            res.status(400).send("Bad Request: quiz has already been created for this topic");
            return;
        }

        // create and save quiz
        const { body } = req.body;
        var newQuiz = {
            topicid: topicid,
            body: body,
            creator: userID,
            createdAt: new Date()
        }
        const query = new Quiz(newQuiz);
        query.save((err, createdQuiz) => {
            if (err) {
                res.status(500).send("Internal server error: " + err.toString());
                return;
            }
            res.setHeader("content-type", "application/json");
            res.status(201).json(createdQuiz)
        });
    } catch (e) {
        res.status(500).send("Internal server error: " + e.toString());
    }
}

// edit exisiting quiz body
const patchQuizHandler = async (req, res, { Topic, Quiz }) => {
    try {
        // check that the user is logged in
        const temp = req.headers["x-user"];
        var userID = temp;
        userID = userID.split("id: ")[1].split("}")[0];
        if (!userID) {
            res.status(403).send("Not authorized: user must be logged in")
        }

        const topicid = req.params.topicID;
        const topic = await Topic.find({_id: topicid}); 
        if (topic.length < 1) {
            res.status(404).send("Not found: Topic does not exist");
            return;
        }

        // get the quiz
        const quiz = await Quiz.find({topicid: topicid});
        if (quiz.length < 1) {
            res.status(400).send("Bad Request: Quiz does not exist");
            return;
        }

        // check that quiz was created by user
        if (quiz.creator != userID) {
            res.status(403).send("Unauthorized: only the quiz creator can make edits");
            return;
        }

        // edit quiz
        const { body } = req.body;
        if (!body) {
            res.status(400).send("Bad Request: must have 'body' text in request body");
            return;
        }
        quiz.body = body;
        quiz.save((err, newQuiz) => {
            if (err) {
                res.status(500).send("Internal server error: " + err.toString());
                return;
            }
            res.setHeader("content-type", "application/json");
            res.status(200).json(newQuiz)
        });
    } catch (e) {
        res.status(500).send("Internal server error: " + e.toString());
    }
}

// delete existing quiz body
const deleteQuizHandler = async (req, res, { Topic, Quiz }) => {
    try {
        // check that the user is logged in
        const temp = req.headers["x-user"];
        var userID = temp;
        userID = userID.split("id: ")[1].split("}")[0];
        if (!userID) {
            res.status(403).send("Not authorized: user must be logged in")
        }

        const topicid = req.params.topicID;
        const topic = await Topic.find({_id: topicid}); 
        if (topic.length < 1) {
            res.status(404).send("Not found: Topic does not exist");
            return;
        }
    
        // get the quiz
        const quiz = await Quiz.find({topicid: topicid});
        if (quiz.length < 1) {
            res.status(400).send("Bad Request: quiz does not exist for this topic");
            return;
        }

        // check that quiz was created by user
        if (quiz.creator != userID) {
            res.status(403).send("Unauthorized: only the quiz creator can delete");
            return;
        }

        Quiz.remove({_id: quiz.id}, (err) => {
            if (err) {
                res.status(500).send("Internal server error: " + err.toString());
                return;
            }
            res.setHeader("content-type", "text/plain").send("Quiz successfully deleted");
        });
    } catch (e) {
        res.status(500).send("Internal server error: " + e.toString());
    }
}

// creates and stores a user's answer to a quiz
// use 'choice' from the request body to get the user's answer
const postAnswerHandler = async (req, res, { Topic, Quiz, Answer }) => {
    try {
        // check that the user is logged in
        const temp = req.headers["x-user"];
        var userID = temp;
        userID = userID.split("id: ")[1].split("}")[0];
        if (!userID) {
            res.status(403).send("Not authorized: user must be logged in")
        }

        const topicid = req.params.topicID;
        const topic = await Topic.find({_id: topicid}); 
        if (topic.length < 1) {
            res.status(404).send("Not found: Topic does not exist");
            return;
        }
    
        // get the quiz
        const quiz = await Quiz.find({topicid: topicid});
        if (quiz.length < 1) {
            res.status(400).send("Bad Request: quiz does not yet exist for this topic");
            return;
        }

        // create answer
        const { choice } = req.body;
        if (choice == undefined) {
            res.status(400).send("Bad Request: request body must contain a T/F choice");
            return;
        }
        const answer = {
            userid: userID,
            quizid: quiz.id,
            answer: choice,
            createdAt: new Date()
        }
        const query = new Answer(answer);
        query.save((err, newAnswer) => {
            if (err) {
                res.status(500).send("Unable to save answer: " + e.toString());
                return;
            }
            res.set("content-type", "application/json");
            res.status(201).json(newAnswer);
        })
    } catch (e) {
        res.status(500).send("Internal server error: " + e.toString());
    }
}



module.exports = {
    getSpecificQuizHandler,
    postQuizHandler,
    patchQuizHandler,
    deleteQuizHandler,
    postAnswerHandler }