PHONY: all
all: servicea serviceb servicec


.PHONY: servicea serviceb servicec
servicea:
	echo ">> Building servicea"
	cd ./servicea/
	go mod tidy
	cd ../
	go build -o ./bin/servicea ./servicea/main.go

serviceb:
	echo ">> Building serviceb"
	go build -o ./bin/serviceb ./serviceb/main.go

servicec:
	echo ">> Building servicec"
	go build -o ./bin/servicec ./servicec/main.go

clean:
	rm ./bin/servicea ./bin/serviceb ./bin/servicec