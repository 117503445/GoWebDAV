#!/bin/bash
# https://github.com/casdoor/casdoor/blob/master/build.sh

# try to connect to google to determine whether user need to use proxy
curl google.com --connect-timeout 1 --max-time 1 -s
if [ $? == 0 ]
then
    echo "Successfully connected to Google, no need to use Go proxy"
    go mod download
else
    echo "Google is blocked, Go proxy is enabled: GOPROXY=https://goproxy.cn,direct"
    GOPROXY=https://goproxy.cn,direct go mod download
fi