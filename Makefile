dev: links.txt
	go run main.go

links.txt:
	cp links.example.txt $@

server: Caddyfile
	sudo caddy start
