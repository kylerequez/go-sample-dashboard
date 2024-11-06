dev:
	@air

templ:
	@templ generate --watch --proxy="http://localhost:1337"

tailwind:
	@npx tailwindcss -i ./src/public/css/input.css -o ./src/public/css/styles.css --watch
