FROM python:3.12-slim

WORKDIR /app

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

RUN python -m grpc_tools.protoc \
    -Iproto \
    --python_out=. \
    --grpc_python_out=. \
    proto/movie.proto

ENV DJANGO_SETTINGS_MODULE=movie_service.settings

CMD ["sh", "-c", "python manage.py migrate && python -u -m movie_service.grpc_server"]
