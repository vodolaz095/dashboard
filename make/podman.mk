podman/up:
	podman-compose up -d
	podman ps

podman/resource:
	podman up -d redis mariadb postgresql influx

podman/down:
	podman-compose down

podman/prune:
	podman system prune -a

podman/build:
	podman build -t dashboard:$(majorVersion).$(minorVersion).$(patchVersion) .
