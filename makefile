dev-up:
	docker compose -f docker-compose.dev.yml up -d --build
dev-down:	
	docker compose -f docker-compose.dev.yml down --volumes
dev-restart:
	docker compose -f docker-compose.dev.yml restart	

test:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down --volumes