all: prepare run clean
run:
	go run main.go
prepare:
	docker run -l mysql=test-by-go -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --rm mysql:5.7
clean:
	docker rm -f $(shell docker ps -f label=mysql=test-by-go -q)
