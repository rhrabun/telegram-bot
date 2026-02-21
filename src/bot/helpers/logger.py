from functools import wraps


def log(func):
    @wraps(func)
    def wrapper(event, context, *args, **kwargs):
        print(f'\n{event}')

    return wrapper
