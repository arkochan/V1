"""
Reusable logging functions for the V1 project scripts.

This module provides standardized logging functions with consistent formatting
and timestamps similar to Go and Next.js applications.
"""

import datetime
import sys
from typing import Optional


# ANSI Color Codes
RESET = "\033[0m"
BOLD = "\033[1m"
DIM = "\033[2m"
GO_COLOR = "\033[38;5;81m"
NEXT_COLOR = "\033[38;5;120m"
SYSTEM_COLOR = "\033[38;5;214m"
ERROR_COLOR = "\033[38;5;196m"
TIME_COLOR = "\033[38;5;244m"


def get_timestamp() -> str:
    """Get the current timestamp in HH:MM:SS format."""
    return datetime.datetime.now().strftime("%H:%M:%S")


def log_system(message: str, bold: bool = True) -> None:
    """Log a system message with timestamp and consistent formatting."""
    timestamp = get_timestamp()
    bold_start = BOLD if bold else ""
    bold_end = RESET if bold else ""

    print(
        f"{TIME_COLOR}{timestamp}{RESET} {SYSTEM_COLOR}[system]{RESET} {bold_start}{message}{RESET}{bold_end}"
    )
    sys.stdout.flush()


def log_error(message: str) -> None:
    """Log an error message with timestamp and consistent formatting."""
    timestamp = get_timestamp()
    print(f"{TIME_COLOR}{timestamp}{RESET} {ERROR_COLOR}[error]{RESET} {message}")
    sys.stdout.flush()


def log_service_output(service_name: str, message: str, color: str) -> None:
    """Log service output with timestamp and service-specific formatting."""
    timestamp = get_timestamp()
    gutter = "[go]  " if service_name.lower() == "go" else "[next]"
    print(f"{TIME_COLOR}{timestamp}{RESET} {color}{gutter}{RESET} {message}")
    sys.stdout.flush()


def log_simple_error(message: str) -> None:
    """Log an error message without timestamp (for compatibility with existing format)."""
    print(f"{ERROR_COLOR}[system]{RESET} {message}")
    sys.stdout.flush()

