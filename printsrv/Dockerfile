FROM python:3.7.6-alpine3.11
RUN apk --update add bash nano
COPY . /srv/printsrv
RUN pip3 install -r /srv/printsrv/requirements.txt
ENTRYPOINT ["python3", "/srv/printsrv/app.py"]
