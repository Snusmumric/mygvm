.PHONY: build
build:
	go build -o bin/mygvm ./

.PHONY: test_install_img
test_install_img: build
	cp bin/mygvm tests/data/mygvm
	sudo docker build -f tests/InstallationTest.Dockerfile -t mygvm_install_test tests/
