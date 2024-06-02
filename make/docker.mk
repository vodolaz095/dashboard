docker/up:
	docker compose up -d
	docker ps

docker/down:
	docker compose down

docker/prune:
	docker system prune -a
