<template>
    <el-container id="layout">
        <el-header>
            <div class="header" style="padding: 18px 0">
                <el-row>
                    <el-col :span="12">
                        <span @click="$router.push('/')" style="cursor: pointer">
                            媒体控制台
                        </span>
                    </el-col>
                    <el-col :span="12">
                        <div style="text-align: right">
                            <!-- <span style="cursor:pointer" @click="handleDisconnect">
                                断开连接
                            </span> -->
                        </div>
                    </el-col>
                </el-row>
            </div>
        </el-header>
        <el-container>
            <el-aside width="auto">
                <!-- <el-menu
                    :collapse="isCollapse"
                    router
                    :default-active="$route.path"
                    class="el-menu-vertical-demo">
                    <el-menu-item v-for="(v,i) in routes" :key="i" :index="v.path">
                        <i class="el-icon-menu"></i>
                        <span slot="title">{{v.name}}</span>
                    </el-menu-item>
                </el-menu> -->
            </el-aside>
            <el-main>
                <nuxt/>
            </el-main>
        </el-container>
    </el-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex"
import OBSWebSocket from "obs-websocket-js"
import { index } from "../mixin"
export default {
    name: "DefaultLayout",
    mixins: [index],
    data(){
        return{
            isCollapse: false,
            routes: [
                {
                    path: '/',
                    name: 'DashBoard',
                    meta: 'global'
                },
                {
                    path: '/obs',
                    name: 'OBS',
                    meta: 'global'
                },
                {
                    path: '/file',
                    name: 'File',
                    meta: 'global'
                },
            ]
        }
    },
    sockets:{
        disconnect(data){
            this.$notify.error({
                title: "socket断开连接",
                message: `${data}:请重新刷新页面,如果依然异常请联系管理员`
            })
        }
    },
    methods:{
        ...mapActions({
            saveConnectInfo: "saveConnectInfo",
        }),
        // async beforeunloadFn(){
        //     await this.$obs.disconnect()
        // },
        // handleDisconnect(){
        //     this.$obs.disconnect()
        //     this.saveConnectInfo(null)
        //     this.$router.push("/connect")
        // }
    },
    destroyed(){
        // window.removeEventListener('beforeunload', e => this.beforeunloadFn(e))
    },
    created(){
        // window.addEventListener('beforeunload', e => this.beforeunloadFn(e))
    },
    mounted(){
        // this.$socket.emit("setIp", "192.168.104.106")
    },
    computed: {
        ...mapGetters({
            connectInfo: "getConnectInfo",
        }),
    }
}
</script>

<style>
html {
  font-family: "Source Sans Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 16px;
  word-spacing: 1px;
  -ms-text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
  -moz-osx-font-smoothing: grayscale;
  -webkit-font-smoothing: antialiased;
  box-sizing: border-box;
}

*, *:before, *:after {
  box-sizing: border-box;
  margin: 0;
}

.button--green {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #3b8070;
  color: #3b8070;
  text-decoration: none;
  padding: 10px 30px;
}

.button--green:hover {
  color: #fff;
  background-color: #3b8070;
}

.button--grey {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #35495e;
  color: #35495e;
  text-decoration: none;
  padding: 10px 30px;
  margin-left: 15px;
}

.button--grey:hover {
  color: #fff;
  background-color: #35495e;
}
#layout{
    height: 100vh;
}
.el-menu{
    height: 100%;
}
.el-header{
    background: url("../assets/images/bg2.jpg");
    background-position-x: -110px;
}
.header{
    /* background: url("../assets/images/bg2.jpg"); */
    height: 60px;
    position: relative;
}
</style>
