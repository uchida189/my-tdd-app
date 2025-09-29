FROM golang:1.24
WORKDIR /app
RUN go install github.com/cucumber/godog/cmd/godog@latest
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
COPY go.mod ./
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]