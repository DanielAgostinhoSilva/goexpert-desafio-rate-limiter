STACK_NAME=desafio-rate-limiter
CONTAINER_APP_NAME=desafio-rate-limiter-app

start:
	docker-compose -p ${STACK_NAME} up -d

build:
	docker-compose -p ${STACK_NAME} up --build -d

stop:
	docker-compose -p ${STACK_NAME} stop

restart: stop start

clean:
	docker-compose -p ${STACK_NAME} down

logs:
	docker-compose -p ${STACK_NAME} logs -f


log-app:
	docker logs -f ${CONTAINER_APP_NAME}

ps:
	docker-compose -p ${STACK_NAME} ps
