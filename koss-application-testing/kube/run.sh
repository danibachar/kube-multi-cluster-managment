#!/bin/bash

exec gunicorn --config ./gunicorn-config.py --log-level debug app:app

#--log-config ./logging.conf 