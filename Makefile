build:
	go build -v -o bin/http cmd/http/*.go

run: build
	bin/http

watch:
	reflex -s -r "\.(go|json|html)$$" --decoration=none make run

log-dev:
	docker logs -f --tail 10 gorilla-ws-development

log-prod:
	docker logs -f --tail 10 gorilla-ws-production

dev:
	bash ./scripts/local.development.sh

prod:
	bash ./scripts/local.production.sh
