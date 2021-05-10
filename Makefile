clean:
	rm -rf build

build: clean
	go build -o build/ste-go src/main/main.go