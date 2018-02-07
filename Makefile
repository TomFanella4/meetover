# REF: Project Arya
# Author: Krutarth Rao (krutarth.rao@hpe.com)


# Network
NET := ranet
SUBNET := 172.20.0.0/24

# Halon
HALON_IMAGE := halon
HALON_CONTAINER := halon
HALON_IP := 172.20.0.2 
halon_download := false
halon_url := http://hpnfiles.rose.rdlabs.hpecorp.net/aruba/pub/cit/LEVEL3_P4/34230/genericx86-p4_essw_cit_10_00_20171117_181048_e8ec0c7.tar.gz

# Backend
UBUNTU_IMAGE := ubuntu:16.10
BACKEND_IMAGE := ra_backend
BACKEND_CONTAINER := backend
BACKEND_IP := 172.20.0.3
BACKEND_GOPATH := ${PWD}
BACKEND_SRC := ${PWD}/ra-backend
TEMP_DIR := ${PWD}/temp/src/ra-backend
TEMP_GOPATH := ${PWD}/temp/
GOOS := linux
GOARCH := amd64

# UI
NODE_IMAGE := node:latest
UI_IMAGE := ui_image
UI_CONTAINER := ui
UI_IP := 172.20.0.4

rm-net:
	-@docker network rm $(NET)
add-net:
	-@docker network create --driver bridge --subnet=$(SUBNET) $(NET)

get-halon:
	@curl --connect-timeout 4 $(halon_url) | docker import - $(HALON_IMAGE)

halon-run:
	@docker run --privileged \
		--detach --net $(NET) \
		--ip $(HALON_IP) \
		--publish 443:443 \
		--name $(HALON_CONTAINER) $(HALON_IMAGE) /sbin/init
	@halon/setup_halon.sh $(HALON_CONTAINER)
halon-restart-rest:
	-@docker exec -it $(HALON_CONTAINER) systemctl restart hpe-restd
halon-clean-container:
	-@docker rm -f $(HALON_CONTAINER)
halon-clean-image:
	-@docker rmi -f $(HALON_IMAGE)

ifeq ($(halon_download), true)
halon-clean: halon-clean-container halon-clean-image get-halon halon-run
halon: halon-clean get-halon halon-run
else
halon-clean: halon-clean-container
halon: halon-clean halon-run
endif


ui-build:
	@docker pull $(NODE_IMAGE)
	@cd ra-web-ui && docker build . --tag $(UI_IMAGE)
ui-run:
	@docker run --detach --net $(NET) --ip $(UI_IP) \
		--publish 5000:5000 --name $(UI_CONTAINER) $(UI_IMAGE)
ui-clean:
	-@docker rm -f $(UI_CONTAINER)
	-@docker rmi $(UI_IMAGE)
	# @cd ra-web-ui # && rm -rf <any dist binaries>

ui: ui-clean ui-build ui-run


# Create binary to send to the container
# Go needs a src direcory and project scructure is different.
# Create a temprary src folder and copy the binary out of it.
backend-dist:
	mkdir -p $(TEMP_DIR)
	cp -r $(BACKEND_SRC)/* $(TEMP_DIR)
	cd $(TEMP_DIR) && GOPATH=$(TEMP_GOPATH) GOBIN=$(TEMP_DIR)/bin go get ./...
	cd $(TEMP_DIR) && GOPATH=$(TEMP_GOPATH) GOOS=$(GOOS) GOARCH=$(GOARCH) go build #
	mv $(TEMP_DIR)/ra-backend $(BACKEND_SRC)/ra-backend
	rm -rf $(TEMP_GOPATH)
	cd $(BACKEND_SRC)

backend-build: backend-dist
	-@docker pull $(UBUNTU_IMAGE)
	-@cd ra-backend && docker build . --tag $(BACKEND_IMAGE)


backend-run:
	-@docker run --detach \
	--env HALON_IP=$(HALON_IP) \
	--net $(NET) \
	--ip $(BACKEND_IP) \
	--name $(BACKEND_CONTAINER) \
	--publish 8080:8080 \
	$(BACKEND_IMAGE)

backend-clean:
	-@docker rm -f $(BACKEND_CONTAINER)
	-@docker rmi -f $(BACKEND_IMAGE)
	-@rm -rf $(TEMP_DIR)


backend: backend-clean backend-build backend-run
clean: ui-clean backend-clean halon-clean rm-net
all: add-net halon backend ui
