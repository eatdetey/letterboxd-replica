FROM python:3.12-slim

WORKDIR /app

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

ENV DJANGO_SETTINGS_MODULE=movie_service.settings
ENV PYTHONPATH=/app:/app/grpc_layer/protobuf

CMD ["sh", "-c", "python manage.py migrate && python -u -m movie_service.grpc_server"]
