start-sqs:
	docker run -p 9324:9324 -p 9325:9325 -v `pwd`/test/queues.conf:/opt/elasticmq.conf softwaremill/elasticmq

test:	test-unit test-integration

test-unit:
	go test ./... --tags=unit -race

test-integration:
	go test ./... --tags=integration -race
