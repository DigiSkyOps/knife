<template>
    <div>
        <el-button @click="$router.push('/')" size="medium" type="text">
            返回
        </el-button>
        <el-input v-model="name" placeholder="文件名称"></el-input>
        <el-tabs type="border-card">
            <el-tab-pane v-for="(v,i) in files" :key="i" :label="v.type">
                <el-table
                    :data="searchFiles(v.data)"
                    max-height="700"
                    style="width: 100%">
                    <el-table-column
                        label="文件名">
                        <template slot-scope="scope">
                            <el-tag @click="downloadFile(scope.row.file,v.type)" style="cursor: pointer">
                                {{scope.row.file}}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column
                        label="修改日期">
                        <template slot-scope="scope">
                            {{dateFormat("YYYY-mm-dd HH:MM:SS",new Date(scope.row.modTime))}}
                        </template>
                    </el-table-column>
                    <el-table-column
                        prop="size"
                        label="文件大小">
                    </el-table-column>
                </el-table>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex"
import GBK from "gbk.js"
export default {
    name: "File",
    async asyncData(context) {
    },
    data(){
        return{
            files: [],
            name: ''
        }
    },
    mounted(){
        this.getHostFile(this.$route.query)
    },
    methods:{
        dateFormat(fmt, date) {
            let ret;
            const opt = {
                "Y+": date.getFullYear().toString(),        // 年
                "m+": (date.getMonth() + 1).toString(),     // 月
                "d+": date.getDate().toString(),            // 日
                "H+": date.getHours().toString(),           // 时
                "M+": date.getMinutes().toString(),         // 分
                "S+": date.getSeconds().toString()          // 秒
                // 有其他格式化字符需求可以继续添加，必须转化成字符串
            };
            for (let k in opt) {
                ret = new RegExp("(" + k + ")").exec(fmt);
                if (ret) {
                    fmt = fmt.replace(ret[1], (ret[1].length == 1) ? (opt[k]) : (opt[k].padStart(ret[1].length, "0")))
                };
            };
            return fmt;
        },
        async getHostFile(query){
            try {
                let filepath = query.filepath.replace(/\\/g,"\/")
                let address = query.address
                let agentport = query.agentport
                let res = await this.$axios.get(`http://${address}:${agentport}/agent/filepath?filepath=${filepath}`)
                let data = res.data

                let types = new Set()
                let arr = []

                data.forEach((v)=>{
                    let type = v.Path.replace(filepath,`${address}`).split("/")[1]
                    types.add(type)
                })
                types.forEach((t)=>{
                    let tmp = []
                    data.forEach((v)=>{
                        let type = v.Path.replace(filepath,`${address}`).split("/")[1]
                        if(type === t){
                            let file = v.Path.split(`${filepath}/${type}/`)[1]
                            tmp.push({
                                file,
                                modTime: v.ModTime,
                                size: this.B2String(v.Size)
                            })
                        }
                    })
                    arr.push({
                        type: t,
                        data: tmp
                    })
                })

                this.files = arr
            } catch (e) {
                this.$notify.error({
                    title: "连接出错",
                    message: "返回首页重新连接",
                })
                this.$router.push("/")
            }
        },
        downloadFile(file,type){
            let path = GBK.URI.encodeURI(`ftp://${this.$route.query.address}/${type}/${file}`)
            window.open(path, "_blank")
        },
        B2String(limit){
            let size = ""
            if(limit < 0.9 * 1024){
                size = limit.toFixed(2) + "B"
            }else if(limit < 0.9 * 1024 * 1024){
                size = (limit/1024).toFixed(2) + "KB"
            }else if(limit < 0.9 * 1024 * 1024 * 1024){
                size = (limit/(1024 * 1024)).toFixed(2) + "MB"
            }else{
                size = (limit/(1024 * 1024 * 1024)).toFixed(2) + "GB"
            }

            let sizeStr = size + ""
            let index = sizeStr.indexOf(".")
            let dou = sizeStr.substr(index + 1 ,2)
            if(dou == "00"){
                return sizeStr.substring(0, index) + sizeStr.substr(index + 3, 2)
            }
            return size
        }
    },
    computed: {
        searchFiles(){
            return (raw)=>{
                let arr = []
                raw.forEach((v)=>{
                    if(v.file.indexOf(this.name) !== -1){
                        arr.push(v)
                    }
                })
                return arr
            }
        }
    },
};
</script>

<style>
.container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family: "Quicksand", "Source Sans Pro", -apple-system, BlinkMacSystemFont,
    "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; /* 1 */
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>
