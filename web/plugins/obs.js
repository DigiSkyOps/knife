import Vue from "vue"
import OBSWebSocket from "obs-websocket-js"

Vue.prototype.$obs = new OBSWebSocket()
