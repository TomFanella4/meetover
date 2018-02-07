# Makefile to deploy meetOver backend to a container

# Backend
UBUNTU_IMAGE := ubuntu:16.10
BACKEND_IMAGE := mo_backend
BACKEND_CONTAINER := backend
# BACKEND_GOPATH := ${PWD}
BACKEND_SRC := ${PWD}/backend
TEMP_DIR := ${PWD}/temp/src/mo-backend
TEMP_GOPATH := ${PWD}/temp/
GOOS := linux
GOARCH := amd64


# Create binary to send to the container
# Go needs a src direcory and project scructure is different.
# Create a temprary src folder and copy the binary out of it.
backend-dist:
	mkdir -p $(TEMP_DIR)
	cp -r $(BACKEND_SRC)/* $(TEMP_DIR)
	cd $(TEMP_DIR) && GOPATH=$(TEMP_GOPATH) GOBIN=$(TEMP_DIR)/bin go get ./... # getting dependencies
	cd $(TEMP_DIR) && GOPATH=$(TEMP_GOPATH) GOOS=$(GOOS) GOARCH=$(GOARCH) go build #build binary
	mv $(TEMP_DIR)/mo-backend $(BACKEND_SRC)/mo-backend # moving binary to backend dir in repo
	rm -rf $(TEMP_GOPATH)
	# cd $(BACKEND_SRC)

backend-build: backend-dist
	-@docker pull $(UBUNTU_IMAGE)
	-@cd backend && docker build . --tag $(BACKEND_IMAGE)


backend-run:
	-@docker run --detach \
#	--env HALON_IP=$(HALON_IP) \
#	--net $(NET) \
#	--ip $(BACKEND_IP) \
	--name $(BACKEND_CONTAINER) \
	--publish 8080:8080 \
	$(BACKEND_IMAGE)

backend-clean:
	-@docker rm -f $(BACKEND_CONTAINER)
	-@docker rmi -f $(BACKEND_IMAGE)
	-@rm -rf $(TEMP_DIR)


backend: backend-clean backend-build
clean: backend-clean
all: backend-clean backend-build backend-run
