BUF ?= buf
GIT ?= git
PYTHON ?= python3

# Список сервисов для генерации/линта.
# Записывать через пробел
SERVICES := user-service movie-service movie-list-service review-service

SERVICE_DIRS := $(addprefix api/,$(SERVICES))

THIRD_PARTY_DIR := third-party
GOOGLEAPIS_DIR := $(THIRD_PARTY_DIR)/googleapis
GRPC_GATEWAY_DIR := $(THIRD_PARTY_DIR)/grpc-gateway
MOVIE_SERVICE_PY_OUT := movie-service/grpc_layer/protobuf
MOVIE_SERVICE_PY_PROTOS := \
	api/movie-service/movie/v1/movie.proto \
	api/movie-list-service/movielist/v1/movielist.proto
REVIEW_SERVICE_PY_OUT := review-service/grpc_layer/protobuf
REVIEW_SERVICE_PY_PROTOS := \
	api/review-service/review/v1/review.proto \
	api/movie-service/movie/v1/movie.proto

.PHONY: all third-party third-party-prepare third-party-github third-party-clean generate lint

all: lint generate

define prepare_local_buf_module
	@rm -f $(1)/buf.lock
	@printf "version: v1\n" > $(1)/buf.yaml
endef

third-party: $(GOOGLEAPIS_DIR)/google/api/annotations.proto $(GRPC_GATEWAY_DIR)/protoc-gen-openapiv2/options/annotations.proto

third-party-prepare:
	$(call prepare_local_buf_module,$(GOOGLEAPIS_DIR))
	$(call prepare_local_buf_module,$(GRPC_GATEWAY_DIR))

$(GOOGLEAPIS_DIR)/google/api/annotations.proto:
	@set -e; \
	mkdir -p $(GOOGLEAPIS_DIR); \
	if ! $(BUF) export buf.build/googleapis/googleapis -o $(GOOGLEAPIS_DIR); then \
		echo "buf export for googleapis failed, fallback to GitHub clone"; \
		rm -rf $(GOOGLEAPIS_DIR); \
		$(GIT) clone --depth 1 https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR); \
	fi
	$(call prepare_local_buf_module,$(GOOGLEAPIS_DIR))

$(GRPC_GATEWAY_DIR)/protoc-gen-openapiv2/options/annotations.proto:
	@set -e; \
	mkdir -p $(GRPC_GATEWAY_DIR); \
	if ! $(BUF) export buf.build/grpc-ecosystem/grpc-gateway -o $(GRPC_GATEWAY_DIR); then \
		echo "buf export for grpc-gateway failed, fallback to GitHub clone"; \
		rm -rf $(GRPC_GATEWAY_DIR); \
		$(GIT) clone --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $(GRPC_GATEWAY_DIR); \
	fi
	$(call prepare_local_buf_module,$(GRPC_GATEWAY_DIR))

third-party-github:
	@rm -rf $(GOOGLEAPIS_DIR) $(GRPC_GATEWAY_DIR)
	@$(GIT) clone --depth 1 https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR)
	@$(GIT) clone --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $(GRPC_GATEWAY_DIR)
	$(call prepare_local_buf_module,$(GOOGLEAPIS_DIR))
	$(call prepare_local_buf_module,$(GRPC_GATEWAY_DIR))

third-party-clean:
	@rm -rf $(GOOGLEAPIS_DIR) $(GRPC_GATEWAY_DIR)
	@mkdir -p $(GOOGLEAPIS_DIR) $(GRPC_GATEWAY_DIR)
	@printf "version: v1\n" > $(GOOGLEAPIS_DIR)/buf.yaml
	@printf "version: v1\n" > $(GRPC_GATEWAY_DIR)/buf.yaml

generate: third-party third-party-prepare
	@set -e; \
	$(PYTHON) -m grpc_tools.protoc \
		-I api/movie-service \
		-I api/movie-list-service \
		--python_out=$(MOVIE_SERVICE_PY_OUT) \
		--grpc_python_out=$(MOVIE_SERVICE_PY_OUT) \
		$(MOVIE_SERVICE_PY_PROTOS); \
	$(PYTHON) -m grpc_tools.protoc \
		-I api/review-service \
		-I api/movie-service \
		--python_out=$(REVIEW_SERVICE_PY_OUT) \
		--grpc_python_out=$(REVIEW_SERVICE_PY_OUT) \
		$(REVIEW_SERVICE_PY_PROTOS); \
	for service in $(SERVICES); do \
		$(BUF) generate api/$$service --template api/$$service/buf.gen.yaml; \
	done

lint: third-party third-party-prepare
	@set -e; \
	for service in $(SERVICES); do \
		$(BUF) lint api/$$service; \
	done
