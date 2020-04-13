const { 
  version : VERSION,
  name : APPNAME
} = require('./package.json')

const ENV = process.env.NODE_ENV || 'development'

/** =================== 全局变量配置 ================= **/
global.ENV = ENV
global.ISPROD = ENV === "production"
global.ISDEV = ENV === "development"
global.VERSION = VERSION
global.APPNAME = APPNAME
/** =================== end ================= **/

//服务入口
require('./src/server')