binaries: | generator loader

generator:
	go build -o generator src/cmd/generator/main.go
	chmod +x generator

loader:
	go build -o loader src/cmd/loader/main.go
	chmod +x loader

clean:
	rm -f generator loader