IMG_NAME = quoteserv
PORT = 4443

.PHONY: run start stop clean tail

run:
	docker run \
	--publish $(PORT):$(PORT) \
	--name $(IMG_NAME) \
	--detach \
	fake-quoteserv

	docker ps

start:
	docker start $(IMG_NAME)

stop:
	docker stop $(IMG_NAME)

clean: stop
	docker rm $(IMG_NAME)

tail:
	docker logs -f $(IMG_NAME)
