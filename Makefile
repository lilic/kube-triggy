build:
	go build -i github.com/lilic/triggy/cmd/triggy

delete:
	kubectl delete deployment nginx-123

.PHONY: build delete
