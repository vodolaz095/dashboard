update_sensor/endpoint:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
	    -H "Content-Type: application/x-www-form-urlencoded" \
		-X POST \
        -d "name=endpoint&value=21" \
		http://localhost:3000/update

update_sensor/redis_subscriber:
	redis-cli publish vodolaz095/dashboard/subscriber `date "+%S"`
