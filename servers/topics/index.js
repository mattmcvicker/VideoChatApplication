"use strict"; //disables several forgiving JavaScript features that make it very easy to introduce subtle bugs
const mysql = require('mysql2');
require('mysql2/promise');
const express = require('express');
//const morgan = require("morgan");
const mongoose = require('mongoose');

const { topicSchema, quizSchema, answerSchema } = require('./schemas');

const { 
    getTopicsHandler, 
    postTopicsHandler } = require('./handlers/topicsHandler');


const { 
    getSpecificTopicsHandler, 
    deleteSpecificTopicsHandler,
    patchSpecificTopicsHandler } = require('./handlers/topicsIDHandler');

const {
    getSpecificQuizHandler,
    postQuizHandler,
    patchQuizHandler,
    deleteQuizHandler,
    postAnswerHandler } = require('./handlers/quizHandler');

const {
    postQueueHandler,
    deleteQueueHandler} = require('./handlers/queuesHandler')

const addr = process.env.ADDR || ":80";

var con = mysql.createConnection({
    host: "441sqldb",
    user: "root",
    password: "mysqlrootpassword",
    database: "wefeuddb"
});

con.connect(function(err) {
    if (err) throw err;
});


const mongoEndpoint = 'mongodb://mongocontainer:27017/mongoDB' // we define the schema inside the code
const [host, port] = addr.split(":")

const Topic = mongoose.model("Topic", topicSchema);
const Quiz = mongoose.model("Quiz", quizSchema);
const Answer = mongoose.model("Answer", answerSchema);

const app = express();
app.use(express.json()); // add JSON request body parsing middleware -- middleware handler function that parses JSON in request body

const connect = () => {
    mongoose.connect(mongoEndpoint)
}
// middleware 
const RequestWrapper = (handler, SchemeAndDbForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDbForwarder)
    }
}
//Topics handlers
app.get("/v1/topics", RequestWrapper(getTopicsHandler, { Topic }));
app.post("/v1/topics", RequestWrapper(postTopicsHandler, { Topic }));

//Topics by ID handlers
app.get("/v1/topics/:topicID", RequestWrapper(getSpecificTopicsHandler, { Topic }));
app.patch("/v1/topics/:topicID", RequestWrapper(patchSpecificTopicsHandler, { Topic, con }));
app.delete("/v1/topics/:topicID", RequestWrapper(deleteSpecificTopicsHandler, { Topic }));

// Queue Handlers
app.post("/v1/queue", RequestWrapper(postQueueHandler, { Topic, con }));
app.delete("/v1/queue", RequestWrapper(deleteQueueHandler, { Topic }));

// quiz handlers
app.route("/v1/topics/:topicID/quiz")
    .get(RequestWrapper(getSpecificQuizHandler, { Topic, Quiz }))
    .post(RequestWrapper(postQuizHandler, { Topic, Quiz }))
    .patch(RequestWrapper(patchQuizHandler, { Topic, Quiz }))
    .delete(RequestWrapper(deleteQuizHandler, { Topic, Quiz }));

// take quiz handler
app.route("/v1/topics/:topicID/quiz/take")
    .post(RequestWrapper(postAnswerHandler, { Topic, Quiz, Answer }))

connect();
mongoose.connection.on('error', console.error)
    .on("disconnected", connect)
    .once('open', main)

async function main() {
    app.listen(port, host, () => { // may need host in the middle?
        console.log(`server listening ${port}`)
    })
}
