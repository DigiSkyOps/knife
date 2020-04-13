const {
    Schema
} = require("mongoose");

const hostSchema = new Schema({
    name: {
        type: String,
        index: true,
        required: true
    },
    ip:{
        type: String,
        index: true,
        required: true
    },
    obsPort: {
        type: String,
        index: true,
        required: true
    },
    agentPort: {
        type: String,
        index: true,
        required: true
    },
    password: {
        type: String,
        index: true,
        required: true
    },
    filePath: {
        type: String,
        index: true,
        required: true
    },
    commands: {
        type: [{
            name: {
                type: String,
                index: true,
            },
            command: {
                type: String,
                index: true,
            },
            commandFilePath: {
                type: String,
                index: true,
            },
            params: {
                type: String,
                index: true,
            },
        }],
        index: true
    },
}, {
    collection: "hosts",
    timestamps: true
})

module.exports = {
    hostSchema,
}
