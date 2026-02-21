FROM python:3.8-slim

WORKDIR /usr/src/app

COPY requirements.txt .
RUN python3 -m pip install --no-cache-dir -r requirements.txt

COPY src/ .

CMD [ "python3", "./main.py" ]
