#!/bin/bash
echo "building the project"
cd ..
go build -o run
cd scripts/ || exit
echo "building the tester"
go build -o tester
cd ..
echo "packaging the required files"
zip -r paxos.zip configs.yaml scripts/tester scripts/init.sh scripts/term.sh run
echo "uploading the binary"
rsync paxos.zip ywi006@uvcluster.cs.uit.no:/home/ywi006/paxos
rm paxos.zip
echo "uploaded successfully"