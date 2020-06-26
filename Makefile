build:
	env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o bin/ruuvi *.go
	# send to raspi
	scp ./bin/ruuvi $$(cat raspi.conf):~/ruuvi