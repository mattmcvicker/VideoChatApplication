"use strict"; //disables several forgiving JavaScript features that make it very easy to introduce subtle bugs
const mysql = require('mysql2');
require('mysql2/promise');
const express = require('express');
const mongoose = require('mongoose');

const { topicSchema, queueSchema } = require('./schemas');

const { 
    getTopicsHandler, 
    postTopicsHandler } = require('./handlers/topicsHandler');


const { 
    getSpecificTopicsHandler, 
    deleteSpecificTopicsHandler,
    patchSpecificTopicsHandler } = require('./handlers/topicsIDHandler');

const {
    postQueueHandler,
    deleteQueueHandler,
    getQueueHandler} = require('./handlers/queuesHandler')

const addr = process.env.ADDR || ":80";



const mongoEndpoint = 'mongodb://mongocontainer:27017/mongoDB' // we define the schema inside the code
const [host, port] = addr.split(":")

const Topic = mongoose.model("Topic", topicSchema);
const Queue = mongoose.model("Queue", queueSchema);

const app = express();
app.use(express.json()); // add JSON request body parsing middleware -- middleware handler function that parses JSON in request body

const connect = () => {
    mongoose.connect(mongoEndpoint, { useNewUrlParser: true, useUnifiedTopology: true })
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
app.patch("/v1/topics/:topicID", RequestWrapper(patchSpecificTopicsHandler, { Topic }));
app.delete("/v1/topics/:topicID", RequestWrapper(deleteSpecificTopicsHandler, { Topic }));

// Queue Handlers
app.post("/v1/queue", RequestWrapper(postQueueHandler, { Queue }));
app.delete("/v1/queue", RequestWrapper(deleteQueueHandler, { Queue }));
app.get("/v1/queue", RequestWrapper(getQueueHandler, { Queue }));

connect();
mongoose.connection.on('error', console.error)
    .on("disconnected", connect)
    .once('open', main)

async function main() {
    app.listen(port, host, () => { // may need host in the middle?
        console.log(`server listening ${port}`)
    })
}
