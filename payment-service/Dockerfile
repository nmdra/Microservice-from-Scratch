FROM python:3.14-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev postgresql-dev

RUN pip install uv

COPY pyproject.toml uv.lock ./

RUN uv sync --no-dev

COPY . .

EXPOSE 8080

CMD ["uv", "run", "python", "main.py"]