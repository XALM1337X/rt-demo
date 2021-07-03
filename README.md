# rt-demo
##To Run the Project: 
    make clean prepare build run 

1) clean tears down the docker containers 
2) prepare spins up the necessary containers 
3) build compiles the go binary inside of the appropriate container 
4) run starts the server process inside of the web container

##Testing
Run `make test` after running the command in the "To Run the Project" section

##Troubleshooting 
###Project failed to run because of new dependency
Example: 
```
dylan@dylan-MS-7C37:~/projects/rt-demo$ make clean prepare build run
docker-compose stop postgres web
Stopping postgres ... done
Stopping web      ... done
docker-compose down 
Removing postgres ... done
Removing web      ... done
Removing network rt-demo_default
docker-compose up -d postgres web
Creating network "rt-demo_default" with the default driver
Creating web      ... done
Creating postgres ... done
docker exec -it web go build -o main cmd/main.go
go: downloading github.com/lib/pq v1.10.2
internal/server/server.go:9:2: no required module provides package github.com/gorilla/mux; to add it:
        go get github.com/gorilla/mux
make: *** [Makefile:12: build] Error 1
```
Solution: 
run `go get github.com/gorilla/mux` and then rerun `make clean prepare build run`

Explanation: 
Running go get will update the go.mod file. Re-running the build process will sync the updated go.mod file to the container and enable it to pull the dependency.