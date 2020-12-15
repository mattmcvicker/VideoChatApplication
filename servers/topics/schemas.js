const Schema = require('mongoose').Schema;

const topicSchema = new Schema({
    name: {type: String, required: true, unique: false},
    votes: {type: Array, required: true, unique: false},
    createdAt: {type: Date, required: true, unique: false},
    creator: {type: Number, required: true, unique: false},
    editedAt: {type: Date, required: false, unique: false}
});

const queueSchema = new Schema({
    userID: {type: Number, required: true, unique: true},
    topicID: {type: String, required: true, unique: false},
    quizAnswer: {type: Boolean, required: true, unique:false},
    roomId: {type: String, required: true, unique: true},
    queueTime: {type: Date, required:false, unique:false}
})

module.exports = { topicSchema, queueSchema }

