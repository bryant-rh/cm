gen-openapi:
	swag init --pd -d ./cmd/server -o ./cmd/server/docs

gen-client:
	swagger generate client -f ./cmd/server/docs/swagger.json -t ./cmd/client

gen-web:
	npx create-react-app web --template typescript


gen-web-client:
	restful-react import --file ./cmd/server/docs/swagger.json  --output ./cmd/web/src/client-bff.ts