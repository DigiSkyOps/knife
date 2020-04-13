const mongoose = require("mongoose");
const debug = require("debug")("connect-db");
const {
    hostSchema,
} = require("./schemas");

const MONGOURL = ISPROD ? process.env.MONGOURL : "mongodb://root:123456@192.168.1.100:40000/knife?authSource=knife"

mongoose.connect(MONGOURL, {
    useNewUrlParser: true,
    useCreateIndex: true,
});
const db = mongoose.connection;

db.on("open", () => {
    debug("数据库连接成功");
});
db.on("error", e => {
    debug(`[error] : 连接失败 ${e}`);
});

const T_HOST = mongoose.model("job", hostSchema)

T_HOST.createIndexes()

module.exports = {
    T_HOST,
}
