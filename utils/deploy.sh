#!/bin/bash

set -x

. utils/deploy.sh.vars

PKG="${APP}_${OS}_${ARCH}"

goxc -arch="$ARCH" -bc="$OS" -os="$OS"
scp $GOPATH/bin/$APP-xc/snapshot/${PKG}.tar.gz "$SERVER:/tmp/" || exit 1
ssh -t $SERVER "tar -C /tmp -vxzf /tmp/${PKG}.tar.gz"
ssh -t $SERVER "sudo supervisorctl stop $SERVICE"
ssh -t $SERVER "cp -r /tmp/$PKG/. $DEST"
ssh -t $SERVER "sudo supervisorctl start $SERVICE"
ssh -t $SERVER "rm -rf /tmp/$PKG/ /tmp/${PKG}.tar.gz"
