#!/bin/bash

REMOTEDEST=/home/debian/go/src/github.com/gdelazzari/udpserial

echo "Creating necessary directories..."
sshpass -p "temppwd" ssh -q -t debian@192.168.1.11 "mkdir -p $REMOTEDEST" > /dev/null

echo "Removing source files from remote..."
sshpass -p "temppwd" ssh -q -t debian@192.168.1.11 "rm -rf $REMOTEDEST/*.go" > /dev/null
echo "Removing web panel built files..."
sshpass -p "temppwd" ssh -q -t debian@192.168.1.11 "rm -rf $REMOTEDEST/panel" > /dev/null

echo "Creating web panel directories..."
sshpass -p "temppwd" ssh -q -t debian@192.168.1.11 "mkdir -p $REMOTEDEST/panel/dist" > /dev/null

echo "Copying source files to remote..."
for file in *.go ; do
  sshpass -p "temppwd" scp -q $file debian@192.168.1.11:$REMOTEDEST/$file > /dev/null
done

echo "Copying web panel built files to remote..."
sshpass -p "temppwd" scp -q -r panel/dist debian@192.168.1.11:$REMOTEDEST/panel/ > /dev/null
