import Vue from 'vue'
import VueSocketIO from 'vue-socket.io'

Vue.use(new VueSocketIO({
    debug: false,
    // connection: "http://192.168.104.106:5000",
    connection: "http://obs.digi-sky.com",
    options: {
        path: "/api/socket"
    }
}))
