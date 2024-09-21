.PHONY: up-success
up-success:
	docker run --name healthd-success --rm -d -p8080:80 -v ./testdata/conf:/conf -v ./testdata/success/scripts:/scripts healthd:latest /healthd

.PHONY: down-success
down-success:
	docker stop healthd-success

.PHONY: up-fail
up-fail:
	docker run --name healthd-fail --rm -d -p8081:80 -v ./testdata/conf:/conf -v ./testdata/fail/scripts:/scripts healthd:latest /healthd

.PHONY: down-fail
down-fail:
	docker stop healthd-fail

.PHONY: up-timeout
up-timeout:
	docker run --name healthd-timeout --rm -d -p8082:80 -v ./testdata/timeout/conf:/conf -v ./testdata/timeout/scripts:/scripts healthd:latest /healthd

.PHONY: down-timeout
down-timeout:
	docker stop healthd-timeout

.PHONY: check-success
check-success:
	curl localhost:8080/healthcheck;echo

.PHONY: check-fail
check-fail:
	curl localhost:8081/healthcheck;echo

.PHONY: check-timeout
check-timeout:
	curl localhost:8082/healthcheck;echo

.PHONY: enter-maintenance
enter-maintenance:
	docker exec healthd-success touch /tmp/maintenance

.PHONY: exit-maintenance
exit-maintenance:
	docker exec healthd-success rm /tmp/maintenance

.PHONY: sighup-healthd
sighup-healthd:
	docker exec healthd-success kill -HUP 1

.PHONY: count-healthd-log
count-healthd-log:
	docker exec healthd-success ls | grep healthd.*log | wc -l | tr -d ' '

.PHONY: cat-healthd-log
cat-healthd-log:
	docker exec healthd-timeout cat healthd.log
