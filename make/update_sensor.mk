export timestamp=$(shell date "+%S")

update_sensor/endpoint:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "name=endpoint&value=$(timestamp)" \
		http://localhost:3000/update

update_sensor/redis_subscriber:
	redis-cli publish vodolaz095/dashboard/subscriber $(timestamp)
