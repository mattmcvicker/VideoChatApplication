"use strict"; //disables several forgiving JavaScript features that make it very easy to introduce subtle bugs
const mysql = require('mysql2');
require('mysql2/promise');
const express = require("express");
//const morgan = require("morgan");
const mongoose = require('mongoose')
const { topicSchema } = require('./schemas')

const { 
    getTopicsHandler, 
    postTopicsHandler } = require('./handlers/topicsHandler')


const { 
    getSpecificTopicsHandler, 
    deleteSpecificTopicsHandler,
    patchSpecificTopicsHandler } = require('./handlers/topicsIDHandler')

const { methodNotAllowedHandler } = require('./handlers/405handler');

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
const Topic = mongoose.model("Topic", topicSchema)
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
app.route("/v1/topics")
    .get(RequestWrapper(getTopicsHandler, { Topic }))
    .post(RequestWrapper(postTopicsHandler, { Topic }))
    .all(methodNotAllowedHandler);

//Topics by ID handlers
app.route("/v1/topics/:topicID")
    .get(RequestWrapper(getSpecificTopicsHandler, { Topic }))
    .patch(RequestWrapper(patchSpecificTopicsHandler, { Topic, con }))
    .delete(RequestWrapper(deleteSpecificTopicsHandler, { Topic }))
    .all(methodNotAllowedHandler);

connect()
mongoose.connection.on('error', console.error)
    .on("disconnected", connect)
    .once('open', main)

async function main() {
    app.listen(port, host, () => { // may need host in the middle?
        console.log(`server listening ${port}`)
    })
}
