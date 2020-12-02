const Schema = require('mongoose').Schema;

const topicSchema = new Schema({
    name: {type: String, required: true, unique: false},
    questions: {type: Array, required: true, unique: false},
    votes: {type: Set, required: true, unique: false},
    createdAt: {type: Date, required: true, unique: false},
    creator: {type: Number, required: true, unique: false},
    editedAt: {type: Date, required: false, unique: false}
})

module.exports = { topicSchema }

