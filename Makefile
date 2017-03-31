startdb:
	@harness/setup.sh start

stopdb:
	@harness/setup.sh stop

push:
	GOOS=linux GOPATH=/Users/fred/data/code/golang go test -c load_test.go
	scp mgo.test root@pste_5.1.0-rc-173-0439.packetsled.com:

push_ricky:
	GOOS=linux GOPATH=/Users/rickyburnett/Packetsled/gopkg/ go test -c load_test.go
	scp mgo.test root@54.245.185.9:
