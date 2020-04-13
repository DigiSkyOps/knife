const express = require('express')
const router = express.Router()
const { T_HOST } = require("../db/server")

router.get('/', async (req, res, next) => {
    const { limit = 10,page = 1 } = req.query
    try{
        // const hosts = await T_HOST.find(null,{ __v: 0 },{ limit: Number(limit), skip: Number(limit * (page - 1)) })
        const hosts = await T_HOST.find(null)
        const total = await T_HOST.countDocuments()
        res.rawData = {
            hosts,
            total
        }
        next()
    }catch(e){
        next(e.message)
    }
})

router.post('/add', async (req, res, next) => {
    const {   name,
            ip,
            obsPort,
            agentPort,
            password,
            filePath,
            commands, } = req.body

    try{
        const count = await T_HOST.countDocuments({ name })
        if(count >= 1){
            next("主机名已存在")
        }else{
            const data = {
                name,
                ip,
                obsPort,
                agentPort,
                password,
                filePath,
                commands,
            }

            await T_HOST.create(data)
            res.rawData = {
                data: "创建成功"
            }
        }
        next()
    }catch(e){
        next(e.message)
    }
});

router.post('/edit', async (req, res, next) => {
    const { id,
            ip,
            obsPort,
            agentPort,
            password,
            filePath,
            commands } = req.body

    try{
        await T_HOST.findByIdAndUpdate(
            {'_id': id},
            {
                ip,
                obsPort,
                agentPort,
                password,
                filePath,
                commands
            }
        )
        res.rawData = "成功"
        next()
    }catch(e){
        next(e.message)
    }
});

router.post('/delete', async(req, res, next)=>{
    const { id } = req.body
    try{
        const count = await T_HOST.countDocuments({
            _id: id
        })
        if(count >= 1){
            await T_HOST.findByIdAndDelete({
                _id: id
            })

            res.rawData = "成功"
            next()
        }else{
            next("主机不存在")
        }
    }catch(e){
        next(e.message)
    }
})

module.exports = router
