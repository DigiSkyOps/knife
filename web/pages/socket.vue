<template>
    <div>
        <el-button type="primary" @click="send">send</el-button>
    </div>
</template>

<script>
export default{
    mounted(){
    },
    sockets: {
        sendMonitor(data){
            console.log(data)
        },
        pingServer(data){
            this.$socket.emit("pongServer", "192.168.104.106")
        },
        pongServer(data){
            console.log(data)
        }
    },
    methods: {
        send(){
            this.$socket.emit("pingServer", "192.168.104.106")
        }
    },
}
</script>
