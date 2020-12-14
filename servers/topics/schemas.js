const Schema = require('mongoose').Schema;

const topicSchema = new Schema({
    name: {type: String, required: true, unique: false},
    questions: {type: Array, required: true, unique: false},
    votes: {type: Set, required: true, unique: false},
    createdAt: {type: Date, required: true, unique: false},
    creator: {type: Number, required: true, unique: false},
    editedAt: {type: Date, required: false, unique: false}
});

const quizSchema = new Schema({
    topicid: {type: Schema.Types.ObjectId, ref: 'Topic', unique: true},
    body: {type: String, required: true, unique: false},
    creator: {type: Number, required: true, unique: false},
    createdAt: {type: Date, required: true},
    editedAt: {type: Date, default: null}
});

const answerSchema = new Schema({
    userid: {type: Number, required: true},
    quizid: {type: Schema.Types.ObjectId, ref: 'Quiz', required: true},
    answer: {type: Boolean, required: true},
    createdAt: {type: Date}
})

const Queue = new Schema({
    userID: {type: Number, required: true, unique: true},
    quizID: {type: Number, required: true, unique: false},
    quizAnswer: {type: Boolean, required: true, unique:false},
    roomId: {type: String, required: true, unique: true},
    queueTime: {type: Date, required:true, unique:false}
})

module.exports = { topicSchema, quizSchema, answerSchema, Queue }

