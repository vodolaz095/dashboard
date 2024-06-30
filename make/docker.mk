docker/up:
	docker compose up -d
	docker ps

docker/resource:
	podman up -d redis mariadb postgresql influx

docker/down:
	docker compose down

docker/prune:
	docker system prune -a
