FROM golang:1.17-rc-buster
WORKDIR /golang/src/rt-demo
RUN go get -u github.com/lib/pq
ENTRYPOINT ["/golang/src/rt-demo/build.sh"]
