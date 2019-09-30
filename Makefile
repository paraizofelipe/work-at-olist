DOCKER_IMAGE=work-at-olist_img

dk-up: ;
	docker-compose up -d;

dk-down: ;
	docker-compose down;

dk-build: ;
	docker build --cache-from $(DOCKER_IMAGE):latest --tag $(DOCKER_IMAGE):$(VERSION) --tag $(DOCKER_IMAGE):latest .;

hk-push: ;
	heroku container:push web --app heroku-work-at-olist

hk-deploy: ;
	heroku container:release web --app heroku-work-at-olist
