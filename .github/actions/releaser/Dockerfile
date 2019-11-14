FROM python:3

RUN pip install python-crontab croniter GitPython

COPY release.py /release.py

ENTRYPOINT ['/usr/local/bin/python', '/release.py']

