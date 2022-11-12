.PHONY: build
build:
	go build -o bin/mygvm ./

.PHONY: test_image
test_image:
	sudo docker build -f tests/Dockerfile -t mygvm_tests ./
