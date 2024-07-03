FROM golang:1.22
WORKDIR "/env-mngr-backend"
COPY . ./
RUN go get .
EXPOSE 3002
CMD [ "go","run","entry.go" ]