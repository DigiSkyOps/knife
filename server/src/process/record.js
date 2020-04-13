const OBSWebSocket = require('obs-websocket-js')
const obs = new OBSWebSocket()

process.on('message', async (m) => {
    try {
        if(m.type === 'data'){
            await obs.connect({ address: m.data.address, password: m.data.password })
            await obs.send("StartRecording")
            let second = m.data.second
            let timer = setInterval(()=>{
                if(second <= 0){
                    clearInterval(timer)
                    obs.send("StopRecording")
                    process.exit(0)
                }else{
                    process.send({
                        type: "time",
                        data: second--
                    })
                }
            },1000)
        }else if(m.type === 'stop'){
            await obs.send("StopRecording")
            process.exit(0)
        }
    } catch (e) {
        process.exit(1)
    }
})
