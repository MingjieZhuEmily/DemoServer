FROM python:3.8-slim-buster

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

WORKDIR /app
COPY . /app

# Add templates and static directories
COPY templates /app/templates
COPY static /app/static
# COPY . .

EXPOSE 8080
CMD ["python", "main.py"]

