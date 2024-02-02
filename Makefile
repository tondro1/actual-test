build:
	@echo "Build CSS"
	npx tailwind -i static/styles.css -o tailwind.css
	go build