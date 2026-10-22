# Production Build for NESCO Local AI Guardian Service
FROM python:3.11-slim

WORKDIR /app

# Install security certificates & curl for container health check
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install
COPY requirements.txt* ./
RUN if [ -f requirements.txt ]; then pip install --no-cache-dir -r requirements.txt; fi

# Copy application source code
COPY app/ ./app/

ENV PYTHONUNBUFFERED=1
ENV PORT=8001

EXPOSE 8001

HEALTHCHECK --interval=15s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8001/healthz || exit 1

ENTRYPOINT ["python", "-m", "app.main"]

