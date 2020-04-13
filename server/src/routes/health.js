const express = require('express')
const router = express.Router()
const isPortReachable = require('is-port-reachable')

router.post("/", async (req, res, next) => {
    const {hosts,port} = req.body
    try {
        let healths = []
        for(let i=0;i<hosts.length;i++){
            const res = await isPortReachable(port, {host: hosts[i]})
            healths.push({
                ip: hosts[i],
                health: res
            })
        }

        res.rawData = healths
        next()
    } catch (e) {
        next(e.message)
    }
})

router.post("/moniter", async (req,res,next) => {
    res.rawData = "ok"
    next()
})


module.exports = router;
