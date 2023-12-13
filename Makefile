# Crazy...
# Recursive wildcards in GNU make
# https://stackoverflow.com/a/18258352
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))
GO_SRC := $(call rwildcard,.,*.go)
OUT_DIR=out
EXE=$(OUT_DIR)/server

dev: links.txt
	go run bin/server/main.go

run: $(EXE)
	./$<

$(EXE): bin/server/main.go $(GO_SRC) | $(OUT_DIR)
	go build -o $@ $<

$(OUT_DIR):
	mkdir -p $@

links.txt:
	cp links.example.txt $@

web: Caddyfile
	sudo caddy start

bg: web run

.PHONY: web run dev bg
