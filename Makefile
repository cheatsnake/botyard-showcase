start:
	docker compose up -d
rebuild:
	docker compose up -d --build
stop:
	docker compose down