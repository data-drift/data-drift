import logging


def get_logger(name):
    logger = logging.getLogger(name)

    logger.setLevel(logging.INFO)

    if not logger.handlers:
        c_handler = logging.StreamHandler()
        c_handler.setLevel(logging.INFO)

        c_format = logging.Formatter("%(name)s - %(levelname)s - %(message)s")
        c_handler.setFormatter(c_format)

        logger.addHandler(c_handler)

    return logger
