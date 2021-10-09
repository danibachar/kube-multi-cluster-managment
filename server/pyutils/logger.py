# -*- coding: utf-8 -*-

from colorlog import ColoredFormatter
import logging
import sys

from utils import EnvVarConfig
cfg = EnvVarConfig()
should_debug_logger = cfg.get('DEBUG_MODE', bool)
# Global vars
# This is a way to make sure we will init the handler only once
# Recomended by stackoverflow
# https://stackoverflow.com/questions/34169547/python-is-there-anyway-to-initialize-a-variable-only-one-time
INIT_LOGGING_HANDLER = True


def get_logger(level=logging.INFO):

    if should_debug_logger:
        level = logging.DEBUG
    global INIT_LOGGING_HANDLER  # Note explanation in top of page
    logger = logging.getLogger()
    logger.setLevel(level)
    if INIT_LOGGING_HANDLER is True:
        handler = get_logging_handler(level)
        logger.addHandler(handler)
        INIT_LOGGING_HANDLER = False
    return logger


def get_logging_handler(level=logging.INFO):
    if should_debug_logger:
        level = logging.DEBUG
    handler = logging.StreamHandler(stream=sys.stdout)
    formatter = ColoredFormatter(
        "%(log_color)s%(levelname)-8s%(reset)s %(asctime)s %(green)s%(name)s"
        "%(reset)s %(message)s",
        reset=True,
        log_colors={
            'DEBUG':    'cyan',
            'INFO':     'blue',
            'WARNING':  'yellow',
            'ERROR':    'red',
            'CRITICAL': 'red,bg_white',
        }
    )
    handler.setFormatter(formatter)
    handler.setLevel(level)
    return handler
