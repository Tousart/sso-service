.PHONY: build start stop restart

build:
	docker compose build

start:
	docker compose up

stop:
	docker compose down

restart: build
	docker compose up