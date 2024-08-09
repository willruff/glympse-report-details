# Copied directly with permission from Glympse, Inc.
import json
import logging
import sys
import time
import random
import traceback


def get_logger() -> logging.Logger:
    """Return a logger configured for json structured logging."""
    # we don't call the _init_logging function because that'd register global
    # loggers, and we don't want globally available loggers, we want loggers
    # which we pass around by reference, since global loggers can end up very
    # cumbersome.
    logger = logging.Logger(str(random.randint(0, 100000)))
    handler = logging.StreamHandler(sys.stdout)
    formatter = _JSONFormatter()
    handler.setFormatter(formatter)
    logger.addHandler(handler)
    # logger = with_fields(logger, app=app, service=service, env=envname)
    return logger


def with_fields(logger, **kwargs) -> logging.Logger:
    """Wrap the provided logger with the supplied keyword arguments."""
    return NestablAdapter(logger, kwargs)


def _init_logging():
    handler = logging.StreamHandler(sys.stdout)
    formatter = _JSONFormatter()
    handler.setFormatter(formatter)
    # logging.basicConfig(level=level, handlers=[handler])
    # replace contents of function with noop so configuration only happens once
    _init_logging.__code__ = (lambda: None).__code__


def _parse_level(raw_level):
    try:
        level = getattr(logging, raw_level.upper())
    except Exception as e:
        raise ValueError(f'Invalid log level: {raw_level}') from e
    return level


class _JSONFormatter(logging.Formatter):
    def __init__(self, *args, json_encoder=None, **kwargs):
        super().__init__(*args, **kwargs)
        self._ignore_fields = set(dir(logging.makeLogRecord({})))
        self.converter = time.gmtime
        self.json_encoder = json_encoder

    def format(self, record):
        t = time.strftime('%Y-%m-%dT%H:%M:%S', self.converter(record.created))
        data = {
            # We can add these back once everything is using proper json logs
            # 'caller': f'{record.pathname}:{record.lineno}',
            # 'function': record.funcName,
            'level': record.levelname.lower(),
            'timestamp': f'{t}.{record.msecs:.0f}Z'
        }

        if isinstance(record.msg, dict):
            data.update(record.msg)
        else:
            data['msg'] = record.getMessage()

        if record.exc_info:
            e_type, e, tb = record.exc_info
            data['error'] = {
              "kind": e_type.__name__,
              "message": str(e),
              "stack": ''.join(traceback.format_exception(*record.exc_info)),
            }

        data.update((attr, getattr(record, attr)) for attr in dir(record)
                    if attr not in self._ignore_fields)

        if self.json_encoder:
            return json.dumps(data, cls=self.json_encoder)
        return json.dumps(data)


class NestablAdapter(logging.LoggerAdapter):
    """LoggerAdapter that merges extra fields if multiple adapters are used."""
    def process(self, msg, kwargs):
        """Override LoggerAdapter implementation to merge extra fields."""
        # order of unpacking matters here. we want kwargs from outer layers
        # to take precedence over the inner layers of adapters
        kwargs['extra'] = {**self.extra, **kwargs.get('extra', {})}
        return msg, kwargs
