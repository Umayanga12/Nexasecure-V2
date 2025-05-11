import logging
import logging.handlers
import os

def setup_logger(log_file="logs/application.log", log_level=logging.DEBUG):
    """
    Set up a sophisticated logging system with console and file handlers.

    Args:
        log_file (str): Path to the log file.
        log_level (int): Logging level (e.g., logging.DEBUG, logging.INFO).

    Returns:
        logging.Logger: Configured logger instance.
    """
    # Create a logger
    logger = logging.getLogger("SecureWalletLogger")
    logger.setLevel(log_level)

    # Avoid duplicate handlers
    if logger.hasHandlers():
        return logger

    # Create a formatter for consistent log messages
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )

    # Console handler for logging to the terminal
    console_handler = logging.StreamHandler()
    console_handler.setLevel(log_level)
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)

    # File handler for logging to a file with rotation
    try:
        os.makedirs(os.path.dirname(log_file), exist_ok=True)
        file_handler = logging.handlers.RotatingFileHandler(
            log_file, maxBytes=5 * 1024 * 1024, backupCount=5  # 5 MB per file, 5 backups
        )
        file_handler.setLevel(log_level)
        file_handler.setFormatter(formatter)
        logger.addHandler(file_handler)
    except (OSError, IOError) as e:
        logger.error(f"Failed to set up file handler: {e}")

    try:
        error_file_handler = logging.FileHandler("logs/error.log")
        error_file_handler.setLevel(logging.ERROR)
        error_file_handler.setFormatter(formatter)
        logger.addHandler(error_file_handler)
    except (OSError, IOError) as e:
        logger.error(f"Failed to set up error file handler: {e}")

    logger.info("Logger initialized successfully.")

    return logger