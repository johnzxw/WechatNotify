#!/usr/bin/env bash
rm -rf /tmp/go*

rm ./main

#��ʼ����go���� -s ȥ�����ű� -w ȥ��DWARF������Ϣ

/usr/local/go/bin/go build -ldflags "-s -w"  echo/main.go

#����upxѹ��
upx -9 main