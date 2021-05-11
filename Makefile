clean:
	rm -rf build/ste-go

build: clean
	go build -o build/ste-go src/main/main.go