.PHONY: build up down restart downv

build:
	docker compose build

up:
	docker compose up

down:
	docker compose down

downv:
	docker compose down -v

restart: build
	docker compose up