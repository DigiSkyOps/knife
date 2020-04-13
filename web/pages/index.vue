<template>
    <div>
        <div style="overflow: hidden">
            <el-button
                type="primary"
                size="mini"
                style="float: right"
                @click="handleAddHost">创建</el-button>
        </div>
        <div class="dashboard-box">
            <el-card class="box-card" v-for="(v,i) in servers" :key="i">
                <div slot="header" class="clearfix">
                    <p style="margin-bottom: 5px">{{v.name}}</p>
                    <el-tag>{{v.ip}}</el-tag>
                    <el-dropdown trigger="click" style="float: right;margin-top: -1px;" @command="(c)=>handleHostCommand(c,v)">
                        <span class="el-dropdown-link">
                            更多<i class="el-icon-arrow-down el-icon--right"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item command="a" icon="el-icon-edit">
                                <span>修改</span>
                            </el-dropdown-item>
                            <el-dropdown-item command="b" icon="el-icon-delete" style="color: red">
                                <span>删除</span>
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                    <img class="loading-img" v-if="v.execing" src="~static/loading.png" alt="" style="float: right;margin-right: 10px;margin-top: 4px;">
                    <el-dropdown v-else trigger="click" style="float: right;margin-right: 10px;margin-top: -1px;" @command="(c)=>exec(c,v)">
                        <span class="el-dropdown-link">
                            启动脚本<i class="el-icon-arrow-down el-icon--right"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item :disabled="!v.agentHealth" v-for="(v1,i) in v.commands" :key="i" :command="v1" icon="el-icon-s-tools">
                                <span>{{v1.name}}</span>
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                    <img class="loading-img" v-if="v.updateAgent" src="~static/loading.png" alt="" style="float: right;margin-right: 10px;margin-top: 4px;">
                    <el-button
                        v-else
                        :disabled="!v.agentHealth"
                        @click="update(v)"
                        style="float: right; padding: 3px 0;margin-right: 10px;"
                        type="text">
                        更新Agent
                    </el-button>
                    <el-button
                        :disabled="!v.agentHealth || !v.ftpHealth"
                        @click="file(v)"
                        style="float: right; padding: 3px 0;margin-right: 10px;"
                        type="text">
                        访问文件
                    </el-button>
                    <el-button
                        :disabled="!v.obsHealth"
                        @click="link(v)"
                        style="float: right; padding: 3px 0;"
                        type="text">
                        连接OBS
                    </el-button>
                </div>
                <div v-for="(v1,i) in healthList" :key="i" class="health" style="padding: 0 15px;">
                    <div v-if="v[v1.key] === null" class="item">
                        {{v1.name}}<el-tag size="mini" style="float: right" type="primary">获取中...</el-tag>
                    </div>
                    <div v-else-if="v[v1.key] === '500'">
                        {{v1.name}}<el-tag size="mini" style="float: right" type="danger">后端服务异常，请检查</el-tag>
                    </div>
                    <div v-else-if="v[v1.key]">
                        {{v1.name}}&nbsp;{{v[v1.version] ? `(ver: ${v[v1.version]})` : ''}}<img class="state-img" src="~assets/images/greenLight.png"/>
                    </div>
                    <div v-else>
                        {{v1.name}}<img class="state-img" src="~assets/images/redLight.png"/>
                    </div>
                </div>
                <el-divider></el-divider>
                <div v-if="v.monitorHealth === null">
                    暂无监控程序启动
                </div>
                <div v-else-if="v.monitorHealth" style="min-height: 85px;max-height: 200px;overflow-y: scroll;padding: 0 10px">
                    <div v-if="v.windowsInfo.length > 0">
                        <div v-for="(v1, i) in v.windowsInfo" :key="i" class="health windows-info">
                            <div v-if="v1[2] === 'Alive'">
                                <span :title="v1[1]" class="monitor-item">
                                    <span class="monitor-window" style="width: 70%;">{{v1[1]}}</span>
                                    <el-button @click="fullWindow(v1[1],v)" style="vertical-align: text-bottom" type="text" size="mini">全屏窗口</el-button>
                                    <el-button @click="closeWindow(v1[1],v)" style="vertical-align: text-bottom" type="text" size="mini">关闭窗口</el-button>
                                    <img class="state-img" src="~assets/images/greenLight.png"/>
                                </span>
                            </div>
                            <div v-else>
                                <span :title="`${v1[1]}(${v1[3]})`" class="monitor-item">
                                    <span class="monitor-window">{{v1[1]}}</span>
                                    &nbsp;
                                    <span class="monitor-die-time">({{v1[3]}})</span>
                                    <img class="state-img" src="~assets/images/redLight.png"/>
                                </span>
                            </div>
                        </div>
                    </div>
                    <div v-else>
                        暂无监控窗口
                    </div>
                </div>
                <div v-else>
                    监控server异常
                </div>
            </el-card>
        </div>
        <el-drawer
            :title="formState === 'add' ? '添加主机' : '修改主机'"
            :before-close="handleClose"
            :visible.sync="showDrawer"
            size="50%">
            <div style="padding: 0 20px">
                <el-form ref="hostForm" :model="form" :rules="rules" style="height: 80vh;overflow-y: scroll;padding-right: 20px">
                    <el-form-item label="主机名称" label-width="150px" prop="name">
                        <el-input :disabled="formState === 'edit'" v-model="form.name" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="主机IP" label-width="150px" prop="ip">
                        <el-input v-model="form.ip" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="主机OBS端口" label-width="150px" prop="obsPort">
                        <el-input v-model="form.obsPort" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="主机OBS密码" label-width="150px" prop="password">
                        <el-input type="password" v-model="form.password" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="主机Agent端口" label-width="150px" prop="agentPort">
                        <el-input v-model="form.agentPort" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="主机FTP路径" label-width="150px" prop="filePath">
                        <el-input v-model="form.filePath" autocomplete="off"></el-input>
                    </el-form-item>
                    <div v-for="(v,i) in form.commands" :key="i" style="overflow: hidden">
                        <el-form-item label="指令名称" label-width="150px">
                            <el-input v-model="v.name" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-form-item label="启动指令" label-width="150px">
                            <el-input v-model="v.command" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-form-item label="启动脚本路径" label-width="150px">
                            <el-input v-model="v.commandFilePath" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-form-item label="启动参数(以逗号隔开)" label-width="150px">
                            <el-input v-model="v.params" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-button type="text" @click="removeCommand(v)" style="float: right">删除</el-button>
                    </div>
                </el-form>
                <div class="drawer-footer">
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button @click="addCommand" type="primary">新增指令</el-button>
                    <el-button v-if="formState === 'add'" type="primary" @click="addHost" :loading="loading">{{ loading ? '提交中 ...' : '确 定' }}</el-button>
                    <el-button v-else type="primary" @click="editHost" :loading="loading">{{ loading ? '提交中 ...' : '确 定' }}</el-button>
                </div>
            </div>
        </el-drawer>
    </div>
</template>

<script>
export default {
    name: "DashBoard",
    data(){
        return{
            servers:[],
            showDrawer: false,
            loading: false,
            form: {
                name: "",
                ip: "",
                obsPort: "",
                agentPort: "",
                password: "",
                filePath: "",
                commands: [{
                    name: "",
                    command: "",
                    commandFilePath: "",
                    params: ""
                }],
                // command: "",
                // commandFilePath: "",
                // params: ""
            },
            rules: {
                name: {required: true, message: '请填写主机名称', trigger: 'change'},
                ip: [
                    {required: true, message: '请填写主机IP', trigger: 'change'},
                    {pattern: /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/, message: '请填写正确主机IP', trigger: 'change'},
                ],
                obsPort: {required: true, message: '请填写OBS端口', trigger: 'change'},
                password: {required: true, message: '请填写OBS密码', trigger: 'change'},
                agentPort: {required: true, message: '请填写Agent端口', trigger: 'change'},
                filePath: {required: true, message: '请填写FTP路径', trigger: 'change'},
            },
            formState: "add",
            id: '',
            webServerHealth: true,
            server: "",
            command: "",
            filePath: "",
            params: "",
            healthList: [{
                name: "FTP服务器",
                key: "ftpHealth"
            },{
                name: "Agent",
                key: "agentHealth",
                version: "agentVersion"
            },{
                name: "Obs",
                key: "obsHealth",
                version: "obsVersion"
            }],
            timer: null,

        }
    },
    mounted(){
        this.getHosts()
    },
    sockets: {
        disconnect(data){
            this.webServerHealth = false
            this.servers.forEach((v1)=>{
                this.$set(v1,"ftpHealth","500")
                this.$set(v1,"agentHealth","500")
                this.$set(v1,"obsHealth","500")
                this.$forceUpdate()
            })
        },
        sendMonitor(data){
            this.servers.forEach((v)=>{
                if(data.ip === v.ip){
                    clearTimeout(v.timer)
                    this.$set(v,"monitorHealth",true)
                    this.$set(v,"windowsInfo",data.data)
                    this.$forceUpdate()
                    v.timer = setTimeout(()=>{
                        this.$set(v,"monitorHealth",false)
                        this.$forceUpdate()
                    },5000)
                }
            })
        },
        pingServer(data){
            if(data.type === 'ftp'){
                data.healths.forEach((v)=>{
                    this.servers.forEach((v1)=>{
                        if(v.ip === v1.ip){
                            this.$set(v1,"ftpHealth",v.health)
                            this.$forceUpdate()
                        }
                    })
                })

            }else if(data.type === 'agent'){
                data.healths.forEach((v)=>{
                    this.servers.forEach((v1)=>{
                        if(v1.ip === v.ip){
                            if(v.health && v1.agentVersion === null){
                                // this.getAgentVersion([{
                                //     ip: v1.ip,
                                //     agentPort: v1.agentPort
                                // }])
                            }
                            this.$set(v1,"agentHealth",v.health)
                            this.$forceUpdate()
                        }
                    })
                })
            }else if(data.type === 'obs'){
                data.healths.forEach((v)=>{
                    this.servers.forEach((v1)=>{
                        if(v1.ip === v.ip){
                            if(v.health && v1.obsVersion === null){
                                // this.getObsVersion([{
                                //     ip: v1.ip,
                                //     obsPort: v1.obsPort,
                                //     password: v1.password
                                // }])
                            }
                            this.$set(v1,"obsHealth",v.health)
                            this.$forceUpdate()
                        }
                    })
                })
            }
        },
        version(data){
            if(data.type === 'agent'){
                data.status.forEach((v)=>{
                    this.servers.forEach((v1)=>{
                        if(v.ip === v1.ip){
                            this.$set(v1,"agentVersion",v.version)
                            this.$forceUpdate()
                        }
                    })
                })
            }else{
                data.status.forEach((v)=>{
                    this.servers.forEach((v1)=>{
                        if(v1.ip === v.ip){
                            this.$set(v1,"obsVersion",v.version ? v.version.obsStudioVersion : null)
                            this.$forceUpdate()
                        }
                    })
                })
            }
        }
    },
    methods:{
        getAgentVersion(hosts){
            let data = {
                hosts: hosts,
                type: "agent",
            }
            this.$socket.emit("version", data)
        },
        getObsVersion(hosts){
            let data = {
                hosts: hosts,
                type: "obs",
            }
            this.$socket.emit("version", data)
        },
        addCommand(){
            this.form.commands.push({
                name: "",
                command: "",
                commandFilePath: "",
                params: ""
            })
        },
        removeCommand(item){
            var index = this.form.commands.indexOf(item)
            if (index !== -1) {
                this.form.commands.splice(index, 1)
            }
        },
        handleAddHost(){
            this.form = {
                name: "",
                ip: "",
                obsPort: "",
                agentPort: "",
                password: "",
                filePath: "",
                commands: [{
                    name: "",
                    command: "",
                    commandFilePath: "",
                    params: ""
                }],
            }
            this.formState = 'add'
            this.showDrawer = true
        },
        async exec(c,v){
            try {
                this.$set(v,"execing",true)
                this.$forceUpdate()
                let data = {
                    command: c.command,
                    filepath: c.commandFilePath,
                    params: [],
                }
                if(c.params && c.params !== ""){
                    data.params = c.params.split(",")
                }
                let res = await this.$axios.post(`http://${v.ip}:5100/agent/shell`,data)
                setTimeout(()=>{
                    this.$notify.success({
                        title: "成功",
                        message: "已执行脚本"
                    })
                    this.$set(v,"execing",false)
                    this.$forceUpdate()
                },2000)
            } catch (e) {
                this.$notify.error({
                    title: "启动失败",
                    message: e.message
                })
                this.$set(v,"execing",false)
                this.$forceUpdate()
            }
        },
        handleHostCommand(c,v){
            if(c === 'a'){
                this.handleEdit(v)
            }else{
                this.deleteHost(v._id)
            }
        },
        handleHealth(hosts){
            this.getFTPHealth(hosts)
            this.getAgentHealth(hosts)
            this.getObsHealth(hosts)
            if(this.webServerHealth){
                setTimeout(()=>{
                    this.handleHealth(hosts)
                },2000)
            }
        },
        async getFTPHealth(hosts){
            let data = {
                type: 'ftp',
                hosts: hosts,
            }
            this.$socket.emit("pingServer",data)
        },
        async getAgentHealth(hosts){
            let data = {
                type: 'agent',
                hosts: hosts,
            }
            this.$socket.emit("pingServer",data)
        },
        async getObsHealth(hosts){
            let data = {
                type: 'obs',
                hosts: hosts,
            }
            this.$socket.emit("pingServer",data)
        },
        async getHosts(){
            try {
                let res = await this.$axios.get("/api/host")
                this.servers = res.data.hosts
                this.servers.forEach((v)=>{
                    v.ftpHealth = null
                    v.agentHealth = null
                    v.agentVersion = null
                    v.obsHealth = null
                    v.obsVersion = null
                    v.monitorHealth = null
                    v.timer = null
                    v.windowsInfo = []
                    v.execing = false
                    v.updateAgent = false
                })

                let healthHosts = []
                let versionHosts = []
                this.servers.forEach((v,i)=>{
                    healthHosts.push({
                        ip: v.ip,
                        ftpPort: '21',
                        agentPort: v.agentPort,
                        obsPort: v.obsPort
                    })
                    versionHosts.push({
                        ip: v.ip,
                        agentPort: v.agentPort,
                        obsPort: v.obsPort,
                        password: v.password
                    })
                })
                this.handleHealth(healthHosts)
                this.getAgentVersion(versionHosts)
                this.getObsVersion(versionHosts)
            } catch (e) {
                this.servers = []
                this.$notify.error({
                    title: "获取失败",
                    message: e.message
                })
            }
        },
        link(v){
            if(v.agentHealth && v.ftpHealth){
                this.$router.push(`/obs?address=${v.ip}&obsport=${v.obsPort}&password=${v.password}&agentport=${v.agentPort}&filepath=${v.filePath}`)
            }else{
                this.$router.push(`/obs?address=${v.ip}&obsport=${v.obsPort}&password=${v.password}`)
            }
        },
        file(v){
            this.$router.push(`/file?address=${v.ip}&agentport=${v.agentPort}&filepath=${v.filePath}`)
        },
        handleEdit(v){
            this.formState = "edit"
            this.id = v._id
            this.form = {
                name: v.name,
                ip: v.ip,
                obsPort: v.obsPort,
                agentPort: v.agentPort,
                password: v.password,
                filePath: v.filePath,
                commands: v.commands,
            }
            this.showDrawer = true
        },
        async addHost(){
            try {
                this.loading = true
                await this.$refs['hostForm'].validate()
                let data = {
                    ...this.form
                }
                await this.$axios.post("/api/host/add",data)
                this.handleClose()
                this.getHosts()
                this.$notify.success({
                    title: "成功",
                    message: "创建主机成功"
                })
            } catch (e) {
                this.$notify.error({
                    title: "创建失败",
                    message: e.message
                })
            } finally{
                this.loading = false
            }
        },
        async editHost(){
            try {
                this.loading = true
                await this.$refs['hostForm'].validate()
                let data = {
                    ...this.form,
                    id: this.id
                }
                await this.$axios.post("/api/host/edit",data)
                this.handleClose()
                this.getHosts()
                this.$notify.success({
                    title: "成功",
                    message: "修改主机成功"
                })
            } catch (e) {
                this.$notify.error({
                    title: "修改失败",
                    message: e.message
                })
            } finally{
                this.loading = false
            }
        },
        async deleteHost(id){
            this.$confirm('此操作将永久删除该主机, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(async () => {
                try {
                    let data = {
                        id
                    }
                    await this.$axios.post("/api/host/delete",data)
                    this.getHosts()
                    this.$notify.success({
                        title: "成功",
                        message: "删除主机成功"
                    })
                } catch (e) {
                    this.$notify.error({
                        title: "失败",
                        message: e.message
                    })
                }
            })

        },
        handleClose(){
            this.$refs['hostForm'].resetFields()
            this.form = {
                name: "",
                ip: "",
                obsPort: "",
                agentPort: "",
                password: "",
                filePath: "",
                commands: [{
                    name: "",
                    command: "",
                    commandFilePath: "",
                    params: ""
                }],
            }
            this.id = ''
            this.formState = 'add'
            this.showDrawer = false
        },
        update(host){
            this.$prompt('请输入地址', '升级Agent', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
            }).then(async ({ value }) => {
                try {
                    if(value && value !== ""){
                        host.updateAgent = true
                        let res = await this.$axios.post(`//${host.ip}:${host.agentPort}/agent/update`,{
                            url: value
                        })
                        setTimeout(()=>{
                            this.$notify.success({
                                title: "成功",
                                message: "成功升级Agent"
                            })
                            this.getAgentVersion([{
                                ip: host.ip,
                                agentPort: host.agentPort,
                                obsPort: host.obsPort,
                                password: host.password
                            }])
                        },10000)
                    }else{
                        this.$notify.info({
                            title: "取消",
                            message: "取消升级"
                        })
                    }
                } catch (e) {
                    this.$notify.error({
                        title: "失败",
                        message: e.message
                    })
                } finally{
                    setTimeout(()=>{
                        host.updateAgent = false
                    },10000)
                }
            })
        },
        async fullWindow(window,host){
            try {
                let data = {
                    window
                }
                let res = await this.$axios.get(`//${host.ip}:${host.agentPort}/agent/fullwindow`,{params: data})
                if(res.data){
                    this.$notify.success({
                        title: "成功",
                        message: `已全屏${window}`
                    })
                }else{
                    this.$notify.error({
                        title: "失败",
                        message: `未全屏${window}`
                    })
                }
            } catch (e) {
                this.$notify.error({
                    title: "失败",
                    message: `未全屏${e.message}`
                })
            }
        },
        async closeWindow(window,host){
            try {
                let data = {
                    window
                }
                let res = await this.$axios.get(`//${host.ip}:${host.agentPort}/agent/closewindow`,{params: data})
                if(res.data){
                    this.$notify.success({
                        title: "成功",
                        message: `已关闭${window}`
                    })
                }else{
                    this.$notify.error({
                        title: "失败",
                        message: `未关闭${window}`
                    })
                }
            } catch (e) {
                this.$notify.error({
                    title: "失败",
                    message: `未关闭${e.message}`
                })
            }
        },
    },
};
</script>

<style lang="less">
.video-item{
    width: 100%;
    height: 500px;
}
.dashboard-box{
    display: flex;
    justify-content: flex-start;
    flex-wrap: wrap;
    .box-card{
        margin: 20px 10px;
        width: calc(50% - 20px);
        .health{
            margin-bottom: 5px;
            .monitor-item{
                width: 100%;
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                .monitor-window{
                    width: 55%;
                    display: inline-block;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }
                .monitor-die-time{
                    color: red;
                    vertical-align: top;
                }
            }
        }
        .state-img{
            width: 15px;
            vertical-align: middle;
            float: right;
            padding: 4px 0;
        }
    }
}
.el-dropdown-link {
    cursor: pointer;
    color: #409EFF;
    vertical-align: middle;
}
.el-icon-arrow-down {
    font-size: 12px;
}
.demonstration {
    display: block;
    color: #8492a6;
    font-size: 14px;
    margin-bottom: 20px;
}
.drawer-footer{
    margin-top: 20px;
    display: flex;
    button{
        flex: 1;
    }
}
.loading-img{
    width: 15px;
    animation: turn 1s linear infinite;
}
@keyframes turn{
    0%{-webkit-transform:rotate(0deg);}
    25%{-webkit-transform:rotate(90deg);}
    50%{-webkit-transform:rotate(180deg);}
    75%{-webkit-transform:rotate(270deg);}
    100%{-webkit-transform:rotate(360deg);}
}
</style>
