.PHONY: clean 
clean:
	docker-compose stop postgres web
	docker-compose down 

.PHONY: prepare 
prepare: 
	docker-compose up -d postgres web

.PHONY: build
build: 
	docker exec -it web go build -o main cmd/main.go

.PHONY: migrate
migrate: 
	# TODO: add migrations here to enforce db schema 
	# TODO: will also need to alias docker command to target the container with the db 

.PHONY: run 
run: 
	docker exec -d web ./main
# TODO: maybe move this to a separate container within the docker network 
.PHONY: test 
test: 
	go test ./tests/...

