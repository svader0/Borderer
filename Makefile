exec_name := Borderer

run: build
	@echo "Running..."
	bin/$(Borderer)

build:
	go build -o bin/$(exec_name) ./main.go