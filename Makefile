DOCKER_IMAGE=work-at-olist_img
HEROKU_APP=olist-call

test: ;
	go test -v -cover ./... -count=1

dk-run: ;
	docker run --name $(HEROKU_APP) -p 80:8989 -e PORT=8989 $(DOCKER_IMAGE)

dk-down: ;
	docker stop $(HEROKU_APP)
	docker rm $(HEROKU_APP)

dk-build: ;
	docker build --cache-from $(DOCKER_IMAGE):latest --tag $(DOCKER_IMAGE):$(VERSION) --tag $(DOCKER_IMAGE):latest .;

hk-push: ;
	heroku container:push web --app $(HEROKU_APP)

hk-deploy: ;
	heroku container:release web --app $(HEROKU_APP)
