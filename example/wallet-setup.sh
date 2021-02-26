#!/usr/bin/env bash

cd "$1"
rm -f venus-wallet

# check port occupancy
port=5678
result=$(lsof -i:${port} | wc -l)
if (($result > 0)); then
  pid=$(lsof -i:5678 | sed -n '2p' | cut -d " " -f2)
  echo "port ${port} processId ${pid} already in use"
  echo "you can kill this processId if it is a wallet process"
  exit 1
fi

# build
go build -o venus-wallet ./cmd
pwd

# wait for successful startup
nohup ./venus-wallet run &

# listen process
result=0
while (($result < 1)); do
  sleep 1s
  result=$(lsof -i:5678 | wc -l)
done

# record pid
pid=$(lsof -i:5678 | sed -n '2p' | cut -d " " -f2)
echo $pid >./example/pid.tmp

# record rpc token
./venus-wallet auth api-info --perm admin >./example/remote-token.tmp
echo "success"
