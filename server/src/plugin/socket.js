let socket = {}
const socketio = require('socket.io')
const _ = require("underscore")
const isPortReachable = require('is-port-reachable')
const request = require('request')
const OBSWebSocket = require('obs-websocket-js')

const hashName = new Map();

//获取io
socket.initSocketio = (server) => {
    io = socketio.listen(server,{
        path: "/api/socket"
    })

    io.on("connection",(socket)=>{
        // socket.on("setIp", (data)=>{
        //     hashName.set(socket.id, data)
        // })
        socket.on("sendMonitor", (data)=>{
            io.emit("sendMonitor",data)
        })

        socket.on("pingServer", async (data)=>{
            let {hosts,type} = data
            let healths = []
            for(let i=0;i<hosts.length;i++){
                const res = await isPortReachable(hosts[i][`${type}Port`], {host: hosts[i].ip})
                healths.push({
                    ip: hosts[i].ip,
                    health: res
                })
            }
            socket.emit("pingServer",{
                type,
                healths
            })
        })

        socket.on("version" , async (data)=>{
            let { hosts,type } = data
            let versions = []
            for(let i=0;i<hosts.length;i++){
                let v = hosts[i]
                const ping = await isPortReachable(v[`${type}Port`], {host: v.ip})
                if(type === 'agent'){
                    if(ping){
                        try {
                            let res = await doGet(`http://${v.ip}:${v[`${type}Port`]}/agent/version`)
                            versions.push({
                                ip: v.ip,
                                version: JSON.parse(res).data
                            })
                        } catch (e) {
                            versions.push({
                                ip: v.ip,
                                version: null
                            })
                        }
                    }else{
                        versions.push({
                            ip: v.ip,
                            version: null
                        })
                    }
                }else if(type === 'obs'){
                    if(ping){
                        try {
                            const obs = new OBSWebSocket()
                            await obs.connect({ address: `${v.ip}:${v[`${type}Port`]}`, password: v.password })
                            let res = await obs.send("GetVersion")
                            versions.push({
                                ip: v.ip,
                                version: res
                            })
                        } catch (e) {
                            console.log(e)
                            versions.push({
                                ip: v.ip,
                                version: null
                            })
                        }
                    }else{
                        versions.push({
                            ip: v.ip,
                            version: null
                        })
                    }
                }
            }

            socket.emit("version",{
                type,
                status: versions
            })
        })

        // socket.on("pingServer", (data)=>{
        //     for(let item of hashName.entries()){
        //         if(item[1] === data){
        //             const targetSocket = _.findWhere(io.sockets.sockets, {id: item[0]})
        //             targetSocket.emit("pingServer", item[0])
        //         }
        //     }
        // })

        // socket.on("pongServer",(data)=>{
        //     io.emit("pongServer", data)
        // })

        // socket.on("disconnect", (reason)=>{
        //     hashName.delete(socket.id)
        // })
    })
}

function doGet(url) {
    return new Promise((resolve, reject) => {
        request.get(url, (error, res, body) => {
            if (!error && res.statusCode == 200) {
                resolve(body);
            } else {
                reject(error);
            }
        })
    })
}

module.exports = socket
