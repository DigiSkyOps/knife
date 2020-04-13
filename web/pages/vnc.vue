<template>
    <div id="vnc">
        <el-button @click="$router.push('/')" size="medium" type="text">
            返回
        </el-button>
        <iframe src="http://192.168.220.70:31034/guacamole/#/client/MQBjAG15c3Fs"></iframe>
    </div>
</template>

<script>
export default{
    name: "vnc"
}
</script>

<style lang="less">
#vnc{
    iframe{
        width: 100%;
        height: calc(100vh - 140px);
    }
}
</style>
