FROM golang:1.17-rc-buster
WORKDIR /golang/src/rt-demo
ENTRYPOINT ["/golang/src/rt-demo/build.sh"]
