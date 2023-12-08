default: all

all: main visualize
	go build .

main: main.go
	go build main.go

visualize: visualize.go
	go build visualize.go

run: 
	go run main.go visualize.go
# clean:
# 	rm -f visualize client