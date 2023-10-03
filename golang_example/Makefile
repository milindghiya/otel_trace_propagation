PHONY: all
all: servicea serviceb servicec


.PHONY: servicea serviceb servicec
servicea:
	echo ">> Building servicea"
	go mod tidy
	go build -o ./bin/servicea ./servicea/main.go

serviceb:
	echo ">> Building serviceb"
	go mod tidy
	go build -o ./bin/serviceb ./serviceb/main.go

servicec:
	echo ">> Building servicec"
	go mod tidy
	go build -o ./bin/servicec ./servicec/main.go

clean:
	rm ./bin/servicea ./bin/serviceb ./bin/servicec