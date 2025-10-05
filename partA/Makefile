
.PHONY: up down logs restart load clean

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=200

restart: down up

load:
	./load-test/run.sh

clean:
	docker compose down -v || true
	docker image rm demoapp:latest || true
