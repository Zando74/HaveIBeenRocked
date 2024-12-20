start-dev:
	docker compose -f ./compose.dev.yml up --build -d db oracle nginx

up-provisioner:
	docker compose -f ./compose.dev.yml up --build provisionner

stop-dev:
	docker compose -f ./compose.dev.yml down