const express = require('express')
const router = express.Router()
const fork = require("child_process").fork

let hashData = new Map()

router.post("/record", async (req,res,next) => {
    try {
        let {address,password,minute} = req.body
        startRecord(address,password,minute)
        res.rawData = "成功"
        next()
    } catch (e) {
        next(e.message)
    }
})

router.post("/stoprecord", async (req,res,next) => {
    try {
        let {address} = req.body
        let obj = hashData.get(address)
        if(obj){
            obj.record.send({
                type: "stop"
            })
            res.rawData = "成功"
        }else{
            res.rawData = -1
        }
        next()
    } catch (e) {
        next(e.message)
    }
})

router.get("/lasttime", async (req,res,next) => {
    try {
        let {address} = req.query
        let obj = hashData.get(address)
        if(obj){
            res.rawData = obj.second
        }else{
            res.rawData = -1
        }
        next()
    } catch (e) {
        next(e.message)
    }
})

async function startRecord(address,password,minute){
    let record = fork(`${__dirname}/../process/record.js`)

    hashData.set(address,{
        second: minute * 60,
        record
    })

    record.send({
        type: "data",
        data: {
            address,
            password,
            second: minute * 60
        }
    })
    record.on("message",(data)=>{
        if(data.type === "time"){
            hashData.set(address,{
                second: data.data,
                record
            })
        }
    })
    record.on('close',(code)=>{
        hashData.delete(address)
    })
    record.on('disconnect',()=>{
        hashData.delete(address)
    })
    record.on('error',(code)=>{
        hashData.delete(address)
    })
    record.on('exit',(code)=>{
        hashData.delete(address)
        record.disconnect()
    })
}

module.exports = router
