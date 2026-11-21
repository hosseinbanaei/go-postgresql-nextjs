tidy:
	cd backend && go mod tidy

back:
	cd backend && air

front:
	cd frontend && npm run dev

.PHONY: tidy back front
