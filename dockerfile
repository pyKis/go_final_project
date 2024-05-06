FROM golang:1.22.0



WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY web ./web
COPY *.go *.db  ./
COPY .env ./ 


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /final_project ./cmd/main.go

CMD ["/final_project"]