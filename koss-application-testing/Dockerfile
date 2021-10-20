FROM python:3.6

RUN mkdir -p /app

COPY /app.py/ /app
COPY /requirements.txt/ /app
RUN pip3 install --no-cache-dir -r /app/requirements.txt

WORKDIR /app
# For local testing
EXPOSE 8081
ENV DEPENDENCIES="{\"destinations\":[{\"target\":\"http://details.default.svc.cluster.local\",\"request_payload_kb_size\":50,\"config\":{\"memory_params\":{\"duration_seconds\":0.2,\"kb_count\":50},\"cpu_params\":{\"duration_seconds\":0.2,\"load\":0.2}}},{\"target\":\"http://reviews.default.svc.cluster.local\",\"request_payload_kb_size\":50,\"config\":{\"memory_params\":{\"duration_seconds\":0.2,\"kb_count\":50},\"cpu_params\":{\"duration_seconds\":0.2,\"load\":0.2}}}]}"
CMD ["python", "app.py"]
