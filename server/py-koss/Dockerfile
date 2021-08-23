FROM python:3.6

RUN mkdir /app
RUN mkdir /app/utils
RUN pip3 install --upgrade pip setuptools wheel
## We copy everything in the root directory
## into our /app directory
COPY /server/py-koss/ /app
COPY /server/pyutils/ /app/utils
RUN pip3 install --no-cache-dir -r /app/requirements.txt
RUN pip3 install --no-cache-dir -r /app/utils/requirements.txt

## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

ENTRYPOINT [] 
