FROM jjanzic/docker-python3-opencv:latest

COPY . /srv/imgproc

RUN pip3 install -r /srv/imgproc/requirements.txt

WORKDIR /srv/imgproc

EXPOSE 8083
ENTRYPOINT ["python3", "/srv/imgproc/app.py"]
