#!/bin/bash

echo "nodejs start...."

#多项目进程启动 最后一个进程用前台输出模式
# cd /data/webapps/pc
# nohup npm start >nohup.out 2>&1 &
# cd /data/webapps/m
# nohup npm start >nohup.out 2>&1

#带node启动方式
# pm2 start build/main.js --no-daemon

#普通nuxt启动方式
yarn start
