bind = ['0.0.0.0:8081']
workers = 4
loglevel = 'info'
worker_class = 'gevent'
worker_connections = 1000
timeout = 5 # Note low timeout to mimik failure
