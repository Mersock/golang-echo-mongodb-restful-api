dev-up:
	docker compose -f docker-compose.dev.yml up -d --build
dev-down:	
	docker compose -f docker-compose.dev.yml down --volumes
dev-restart:
	docker compose -f docker-compose.dev.yml restart	

prod-up:
	docker compose -f docker-compose.prod.yml up -d --build
prod-down:	
	docker compose -f docker-compose.prod.yml down --volumes	

test:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down --volumes