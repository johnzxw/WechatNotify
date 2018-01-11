#!/usr/bin/env bash
rm -rf /tmp/go*

rm ./main

#开始编译go程序 -s 去掉符号表 -w 去掉DWARF调试信息

/usr/local/go/bin/go build -ldflags "-s -w"  echo/main.go

#开启upx压缩
upx -9 main