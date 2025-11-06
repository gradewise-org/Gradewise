FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app


# Python deps
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# App
COPY . .

# Runtime
ENV PORT=8080 \
    FLASK_ENV=production
EXPOSE 8080

# Use gunicorn in cluster
CMD ["gunicorn", "-w", "2", "-b", "0.0.0.0:8080", "app:app"]