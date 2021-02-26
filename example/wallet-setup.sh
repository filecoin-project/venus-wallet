#!/usr/bin/env bash

cd "$1"
rm -f venus-wallet
port=5678
result=$(lsof -i:${port} | wc -l)
if (($result > 0)); then
  pid=$(lsof -i:5678 | sed -n '2p' | cut -d " " -f2)
  echo "port ${port} processId ${pid} already in use"
  echo "you can kill this processId if it is wallet process"
  exit 1
fi

go build -o venus-wallet ./cmd
pwd
nohup ./venus-wallet run &
result=0
while (($result < 1)); do
  sleep 1s
  result=$(lsof -i:5678 | wc -l)
done
pid=$(lsof -i:5678 | sed -n '2p' | cut -d " " -f2)
echo $pid >./example/pid.tmp
./venus-wallet auth api-info --perm admin  >./example/remote-token.tmp
echo "success"
