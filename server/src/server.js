const express = require('express')
const logger = require('morgan')
const cookieParser = require('cookie-parser')
const bodyParser = require('body-parser')
const qs = require('qs')
const app = express()
const rp = require('request-promise')
const socket = require("./plugin/socket")
// const {
//     T_USER
// } = require("./db/server")

app.use(logger('dev'))
app.use(bodyParser.json({limit: '50mb'}))
app.use(bodyParser.urlencoded({ limit: '50mb',extended: false }))
app.use(cookieParser())

// 认证
// app.use('/', async (req,res,next)=>{
//     let token = req.cookies.token
//     let host = req.hostname
//     let options = {
//         method: 'POST',
//         uri: 'https://sso.digisky.com/api/user/userAuth',
//         form:{
//             host: ISPROD ? host : 'digi-sky.com:8000',
//             token: token
//         },
//         headers: {
//             'Content-Type':'application/x-www-form-urlencoded'
//         }
//     }
//     try{
//         let data = await rp(options)
//         data = JSON.parse(data)
//         if(data.msg === '成功 ' && data.code === 1){
//             const user = await T_USER.find({
//                 userName: data.data.user.username,
//                 nick: data.data.user.nick,
//             }, {
//                 __v: 0,
//                 createdAt: 0,
//                 updatedAt: 0,
//             })
//             if(user.length >= 1){
//                 if(JSON.stringify(req.body) != "{}"){
//                     req.body = qs.parse(req.body)
//                 }
//                 next()
//             }else{
//                 res.status(200).send({
//                     code: 401,
//                     data: '无此用户'
//                 })
//             }
//         }else{
//             res.status(200).send({
//                 code: 401,
//                 data: '未登录'
//             })
//         }
//     }catch(e){
//         res.status(200).send({
//             code: 401,
//             data: e.message
//         })
//     }
// })

// 路由
app.use('/api/host', require('./routes/index'))
app.use('/api/health', require('./routes/health'))
app.use('/api/media', require('./routes/media'))

// 统一处理响应数据handler
app.use(require("./middleware/apiHandler"))

// catch 404 and forward to error handler
app.use(require("./middleware/404Handler"))

// error handler
app.use(require("./middleware/errorHandler"))

const serverInfo = `
| ******************************************************************************** |
| *******  server is running at 0.0.0.0:5000 && [mode]: ${ENV} ******************* |
| *******  version: ${VERSION}     &      app name: ${APPNAME} ******************* |
| ******************************************************************************** |
`;
const http = require("http").createServer(app)
socket.initSocketio(http)

http.listen("5000", () => console.log(serverInfo));
