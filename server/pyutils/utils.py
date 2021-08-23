# -*- coding: utf-8 -*-

from datetime import datetime
import os
import numpy as np
from flask.json import JSONEncoder


class Tools(object):

    @staticmethod
    def chunks(l, n):
        """Yield successive n-sized chunks from l."""
        for i in range(0, len(l), n):
            yield l[i:i + n]


class EnvVarConfig(object):
    def __init__(self):
        pass

    def get(self, var_name, var_type):
        var = os.environ[var_name]  # Raises KeyError if key doesn't exist
        if var_type is str:
            return var
        elif var_type is int:
            return int(var)
        elif var_type is list:
            return var.split(',')
        elif var_type is bool:
            return var.lower() == 'true'
        else:
            return None


class CustomJSONEncoder(JSONEncoder):
    """Minify JSON output and serialize datetime as ISO strings instead of RFC 822."""

    item_separator = ','
    key_separator = ':'

    def default(self, obj):
        try:
            if isinstance(obj, datetime):
                return obj.isoformat()
            elif isinstance(obj, np.integer):
                return int(obj)
            elif isinstance(obj, np.floating):
                return float(obj)
            elif isinstance(obj, np.ndarray):
                return obj.tolist()
            elif hasattr(object, 'to_dict'):
                print("CustomJSONEncoder, has to_dict attr")
            iterable = iter(obj)
        except TypeError:
            pass
        else:
            return list(iterable)
        return JSONEncoder.default(self, obj)


class MicroserviceError(Exception):
    status_code = 400

    def __init__(self, message, status_code=None, payload=None):
        # type: (str, int, dict) -> None

        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload

    def to_dict(self):
        # type: () -> dict
        rv = dict(self.payload or ())
        rv['message'] = self.message
        return rv
