#!/bin/sh -ex
#########################################################################
# File Name: prepare-dirs.sh
# Author: nian
# Blog: https://whoisnian.com
# Mail: zhuchangbao1998@gmail.com
# Created Time: 2022年06月01日 星期一 19时56分12秒
#########################################################################

mkdir_if_not_exist() {
    if [ ! -d "$1" ]; then
        mkdir "$1" # && chown 500:500 "$1"
    fi
}

mkdir_if_not_exist "/public/assets/"
