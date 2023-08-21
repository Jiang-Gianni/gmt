qtc:
	qtc -dir markdown -ext html

build:
	go build -ldflags "-w -s" -gcflags=all="-l -B" gmt.go