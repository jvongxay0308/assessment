test:
	docker compose up -d pgsql15
	docker compose up go
	docker compose down --volumes 

