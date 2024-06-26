export timestamp=$(shell date "+%S")

update_sensor/endpoint/update:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "name=endpoint&value=$(timestamp)" \
		http://localhost:3000/update

update_sensor/endpoint/increment:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "name=endpoint&value=100" \
		http://localhost:3000/increment

update_sensor/endpoint/decrement:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "name=endpoint&value=50" \
		http://localhost:3000/decrement

update_sensor/redis_subscriber:
	redis-cli publish vodolaz095/dashboard/subscriber $(timestamp)
