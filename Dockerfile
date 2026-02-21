FROM python:3.10.5-slim

WORKDIR /usr/src/app

RUN python3 -m pip install -U pip

COPY requirements.txt .
RUN python3 -m pip install --no-cache-dir -r requirements.txt

COPY src/ .

CMD [ "python3", "-u", "./main.py" ]
