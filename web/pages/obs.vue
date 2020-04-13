<template>
    <div>
        <el-button @click="$router.push('/')" size="medium" type="text">
            返回
        </el-button>
        <div id="info" v-if="isConnect">
            <!-- <div id="sceneScene">
                <el-button @click="getSceneList">获取场景列表</el-button>
                <span v-if="scene && scene.scenes.length > 0">
                    <el-tag v-for="(v,i) in scene.scenes" :key="i">
                        {{v.name}}
                    </el-tag>
                </span>
            </div>
            <div id="sourceScene">
                <el-button @click="getSourcesList">获取资源列表</el-button>
                <el-select v-model="source" placeholder="请选择" @change="getSourceSettings">
                    <el-option
                        v-for="item in sources"
                        :key="item.name"
                        :label="item.name"
                        :value="item.name">
                    </el-option>
                </el-select>
                <div>{{sourceSettings}}</div>
            </div> -->
            <el-card class="box-card">
                <div slot="header" class="clearfix">
                    <el-button @click="getStreamInfo" v-if="!initing">获取串流信息</el-button>
                    <el-button v-if="initing" :loading="initing">初始化串流中...</el-button>
                    <el-button
                        style="float: right; padding: 3px 0"
                        type="text"
                        :disabled="stream.server === ''"
                        @click="toggleStream">
                        {{!streamStatus ? "开始串流" : "停止串流"}}
                    </el-button>
                    <el-button
                        style="float: right; padding: 3px 0"
                        type="text"
                        :disabled="stream.server === ''" @click="handleRecord">
                        {{!recordStatus ? "开始录像" : timeTip !== -1 ? `剩余:${timeTip}s` : '停止录像'}}
                    </el-button>
                    <el-button
                        v-if="$route.query.agentport && $route.query.filepath"
                        type="text"
                        @click="goFile"
                        style="float: right; padding: 3px 0">
                        查看文件
                    </el-button>
                </div>
                <div style="overflow: hidden">
                    <el-tag>Address：{{stream.server}}</el-tag>
                    <el-tag>Key：{{stream.key}}</el-tag>
                    <el-tag v-if="lastFile" style="float:right;cursor: pointer" @click="downloadFile">{{lastFile}}</el-tag>
                </div>
            </el-card>
            <div id="videoScene">
                <video
                    ref="video"
                    class="video-js vjs-default-skin video-item"
                    controls
                    style="margin-bottom: 20px"></video>
                <el-button @click="loadVideo" :disabled="!streamStatus">开始播放</el-button>
                <el-button @click="stopVideo" :disabled="!streamStatus">停止播放</el-button>
                <el-button @click="reloadVideo" :disabled="!streamStatus">重新加载</el-button>
            </div>
        </div>
        <el-dialog :close-on-click-modal="false" center title="录制配置" :visible.sync="recordSettingVisible">
            <el-form :model="form">
                <el-form-item label="录制模式" label-width="120px">
                    <el-select style="width: 100%" v-model="form.type" placeholder="请选择录制模式">
                        <el-option label="定时录制" value="time"></el-option>
                        <el-option label="手动录制" value="manual"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item v-if="form.type === 'time'" label="录制时间(分钟)" label-width="120px">
                    <el-input type="number" v-model="form.time" autocomplete="off"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="recordSettingVisible = false">取消</el-button>
                <el-button :disabled="sending" type="primary" @click="handleRecordType">{{sending ? "请求中..." : "录制"}}</el-button>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex"
import OBSWebSocket from "obs-websocket-js"
import { sha256 } from "js-sha256"
import { Base64 } from 'js-base64'
import { v4 as uuidv4 } from 'uuid'
import { index } from "../mixin"
import GBK from "gbk.js"
import config from "../config"
export default {
    name: "Obs",
    mixins: [ index ],
    data(){
        return{
            recordSettingVisible: false,
            server: "",
            isConnect: false,
            scene: null,
            sources: [],
            source: '',
            sourceSettings: null,
            streamSettings: null,
            streamStatus: false,
            recordStatus: false,
            stream:{
                server: '',
                key: ''
            },
            player: null,
            timeTip: -1,
            timer: null,
            form: {
                type: '',
                time: 1
            },
            sending: false,
            initing: false,
            lastFile: ''
        }
    },
    mounted(){
        this.connectObs(this.$route.query)
    },
    watch:{
        isConnect(newV,oldV){
            this.$nextTick(()=>{
                if(newV){
                    this.player = this.videojs(this.$refs.video)
                }else if(newV){
                    this.player.reset()
                    this.player = null
                }
            })
        }
    },
    methods:{
        goFile(){
            let query = this.$route.query
            this.$router.push(`/file?address=${query.address}&agentport=${query.agentport}&filepath=${query.filepath}`)
        },
        downloadFile(file,type){
            let path = GBK.URI.encodeURI(`ftp://${this.$route.query.address}/video/${this.lastFile}`)
            window.open(path, "_blank")
        },
        async getLastFile(){
            try {
                let query = this.$route.query
                let res = await this.$axios.get(`http://${query.address}:${query.agentport}/agent/lastfile?filepath=${query.filepath}/video`)
                this.lastFile = res.data.Path.split("/")[res.data.Path.split("/").length-1]
            } catch (e) {
                this.$notify.error({
                    title: "获取失败",
                    message: "获取最新文件失败，请重试",
                })
            }
        },
        async connectObs(query){
            try {
                if(this.isConnect){
                    await this.$obs.disconnect()
                    this.resetInfo()
                }

                let data = {
                    address: `${query.address}:${query.obsport}`,
                    password: query.password
                }
                await this.$obs.connect(data)
                this.initListen(this.$obs)
                this.isConnect = true
            } catch (e) {
                this.$notify.error({
                    title: "连接出错",
                    message: "返回首页重新连接",
                })
                this.$router.push("/")
            }
        },
        resetInfo(){
            this.isConnect = false
            this.player = null
            this.scene = null
            this.sources =  []
            this.source = ''
            this.sourceSettings = null
            this.streamSettings = null
            this.streamStatus = false
            this.recordStatus = false
            this.stream = {
                server: '',
                key: ''
            }
        },
        initListen(obs){
            obs.on("ConnectionOpened",()=>{
            })
            obs.on('ConnectionClosed',()=>{

            })
            obs.on('error', err => {
                Notification.error({
                    title: "socket error!",
                    message: err
                })
            })
            obs.on('RecordingStarted', async ()=>{
                if(this.timer === null && this.timeTip === -1){
                    let lasttime = await this.$axios.get("/api/media/lasttime",{params: {address: this.$route.query.address}})
                    this.timeTip = lasttime.data
                    this.timer = setInterval(() => {
                        if(this.timeTip <= 0){
                            clearInterval(this.timer)
                            this.timer = null
                        }else{
                            this.timeTip--
                        }
                    }, 1000)
                }
                this.recordStatus = true
            })
            obs.on('RecordingStopping', ()=>{
                clearInterval(this.timer)
                this.timer = null
                this.timeTip = -1
                this.recordStatus = false
                this.getLastFile()
            })
        },
        async getSceneList(){
            let data = await this.$obs.send('GetSceneList')
            this.scene = data
        },
        async getSourcesList(){
            let data = await this.$obs.send('GetSourcesList')
            this.sources = data.sources
        },
        async getSourceSettings(){
            let data = await this.$obs.send('GetSourceSettings',{
                sourceName: this.source
            })
            this.sourceSettings = data
        },
        async getStreamInfo(){
            let res = await this.$obs.send("GetStreamingStatus")
            this.streamStatus = res.streaming
            this.recordStatus = res.recording

            let data = await this.$obs.send('GetStreamSettings')
            this.streamSettings = data
            if(data.settings.server === "" && data.settings.key === ""){
                this.initing = true
                let key = uuidv4()
                await this.$obs.send("SetStreamSettings",{
                    settings: {
                        server: `${config.rtmpServer}/hls`,
                        key: key
                    }
                })
                setTimeout(()=>{
                    this.initing = false
                    this.stream.server = `${config.rtmpServer}/hls`
                    this.stream.key = key
                }, 1000)
            }else{
                this.stream.server = data.settings.server
                this.stream.key = data.settings.key
                if(res.streaming && data.settings.server !== '' && data.settings.key !== ''){
                    this.player.src({
                        src: `${config.videoServer}/hls/${data.settings.key}/index.m3u8`,
                        type: 'application/x-mpegURL'
                    })
                    this.player.play()
                }
            }

            let lasttime = await this.$axios.get("/api/media/lasttime",{params: {address: this.$route.query.address}})
            if(lasttime.data !== -1){
                this.timeTip = lasttime.data
                this.timer = setInterval(() => {
                    if(this.timeTip <= 0){
                        clearInterval(this.timer)
                        this.timer = null
                    }else{
                        this.timeTip--
                    }
                }, 1000)
            }
        },
        async toggleStream(){
            let res = await this.$obs.send("StartStopStreaming")
            if(res.status === 'ok'){
                this.streamStatus = !this.streamStatus
                if(!this.streamStatus){
                    this.player.reset()
                }else{
                    this.loadVideo()
                }
            }
        },
        handleRecord(){
            if(!this.recordStatus && this.timeTip === -1){
                this.recordSettingVisible = true
            }else if(this.recordStatus && this.timeTip === -1){
                this.$confirm('是否停止录制？', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.stopRecord()
                })
            }else if(this.recordStatus && this.timeTip !== -1){
                this.$confirm('是否强制停止定时录制？', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.stopTimeRecord()
                })
            }
        },
        async handleRecordType(){
            this.sending = true
            if(this.form.type === 'time'){
                this.startTimeRecord(this.form.time)
                this.recordSettingVisible = false
            }else{
                this.startRecord()
                this.recordSettingVisible = false
            }
            this.sending = false
        },
        async startRecord(){
            await this.$obs.send("StartRecording")
        },
        async stopRecord(){
            await this.$obs.send("StopRecording")
        },
        async startTimeRecord(minute){
            try {
                let data = {
                    address: `${this.$route.query.address}:${this.$route.query.obsport}`,
                    password: this.$route.query.password,
                    minute: minute
                }
                await this.$axios.post("/api/media/record", data)
                this.timeTip = minute * 60
                this.timer = setInterval(() => {
                    if(this.timeTip <= 0){
                        clearInterval(this.timer)
                        this.timer = null
                    }else{
                        this.timeTip--
                    }
                }, 1000)
            } catch (e) {
                this.$notify.error({
                    title: "失败",
                    message: `定时录制失败: ${e}`
                })
            }
        },
        async stopTimeRecord(){
            try {
                let data = {
                    address: `${this.$route.query.address}:${this.$route.query.obsport}`,
                }
                await this.$axios.post("/api/media/stoprecord", data)
                clearInterval(this.timer)
                this.timer = null
                this.timeTip = -1
            } catch (e) {
                this.$notify.error({
                    title: "失败",
                    message: `停止定时录制失败: ${e}`
                })
            }
        },
        loadVideo(){
            this.player.src({
                src: `http://192.168.220.70:30186/hls/${this.stream.key}/index.m3u8`,
                type: 'application/x-mpegURL'
            })
            this.player.play()
        },
        stopVideo(){
            this.player.reset()
        },
        reloadVideo(){
            this.player.reset()
            this.player.src({
                src: `http://192.168.220.70:30186/hls/${this.stream.key}/index.m3u8`,
                type: 'application/x-mpegURL'
            })
            this.player.play()
        }
    },
    computed: {
        ...mapGetters({
            connectState: "getConnectState"
        }),
    },
};
</script>

<style>
.video-item{
    width: 100%;
    height: 500px;
}
#videoScene{
    margin-top: 20px;
}
</style>
