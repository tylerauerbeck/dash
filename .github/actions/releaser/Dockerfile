FROM python:3

RUN pip install python-crontab croniter GitPython

COPY entrypoint.sh /entrypoint.sh
COPY release.py /release.py

ENTRYPOINT ['/entrypoint.sh']

