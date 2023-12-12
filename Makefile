dev: links.txt
	go run main.go

run: startup
	./$<

startup: main.go
	go build -o $@ $<

links.txt:
	cp links.example.txt $@

server: Caddyfile
	sudo caddy start

.PHONY: server run dev
