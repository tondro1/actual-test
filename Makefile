build:
	@echo "Build CSS"
	npx tailwind -i static/styles.css -o static/tailwind.css
	@echo "Build Go project"
	go build
	./actual-test

watch:
	@echo "Watch CSS"
	npx tailwind --minify --watch -i static/styles.css -o static/tailwind.css