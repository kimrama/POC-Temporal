up:
	docker compose up --build

down:
	docker compose down -v

logs:
	docker compose logs -f backend worker-k8s worker-cert

start:
	curl -X POST http://localhost:8000/workflows \
		-H 'Content-Type: application/json' \
		-d '{"app_name":"demo-app","namespace":"demo-ns","cluster":"dev"}'
