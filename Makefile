hello world mine:
	echo $@

greeting: world
	echo hi



.PHONY: bash go gofmt golangci-lint
bash go gofmt golangci-lint:
	@docker build \
		--tag $@ \
		--build-arg "user_id=$(shell id -u)" \
		--build-arg "group_id=$(shell id -g)" \
		--build-arg "home=${HOME}" \
		--build-arg "workdir=${PWD}" \
		--target $@ . \
		>/dev/null