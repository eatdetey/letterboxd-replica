FROM python:3.12-slim

WORKDIR /app

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

RUN mkdir -p grpc_layer/protobuf && \
    python -m grpc_tools.protoc \
        -I=proto \
        --python_out=grpc_layer/protobuf \
        --grpc_python_out=grpc_layer/protobuf \
        proto/review.proto

ENV DJANGO_SETTINGS_MODULE=review_service.settings
ENV PYTHONPATH=/app:/app/grpc_layer/protobuf

CMD ["sh", "-c", "python manage.py migrate && python -u -m review_service.grpc_server"]
