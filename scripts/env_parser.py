"""
Environment variable parser module.
Parses .env files with variable expansion support.
"""

import os
import re
from typing import Dict


# ANSI Color Codes
RESET = "\033[0m"
SYSTEM_COLOR = "\033[38;5;214m"
ERROR_COLOR = "\033[38;5;196m"


# ============= Environment Parser =============
VAR_PATTERN = re.compile(r"(?<!\\)\$\{([^}]+)\}")  # matches ${VAR}, not \${VAR}


def expand_vars(value: str, env: Dict[str, str]) -> str:
    """Expand ${VAR} using values from env."""
    def repl(match):
        var = match.group(1)
        return env.get(var, "")

    value = VAR_PATTERN.sub(repl, value)
    return value.replace(r"\${", "${")  # unescape \${VAR}


def parse_env_line(line: str):
    """Parse a single .env line into (key, value) or None."""
    line = line.strip()

    # Skip empty or full-line comments
    if not line or line.startswith("#"):
        return None

    # Remove 'export '
    if line.startswith("export "):
        line = line[len("export ") :].lstrip()

    # Split key=value
    if "=" not in line:
        return None

    key, value = line.split("=", 1)

    key = key.strip()
    value = value.strip()

    # Remove inline comments (only if outside quotes)
    if value and not (value.startswith('"') or value.startswith("'")):
        if " #" in value:
            value = value.split(" #", 1)[0].rstrip()
        elif "#" in value:
            value = value.split("#", 1)[0].rstrip()

    # Remove surrounding quotes
    if (value.startswith('"') and value.endswith('"')) or (
        value.startswith("'") and value.endswith("'")
    ):
        value = value[1:-1]

    return key, value


def load_dotenv(path: str = ".env") -> Dict[str, str]:
    """
    Parse .env file with variable expansion.
    Returns dict of variables defined in the file.
    """
    env = {}
    base = os.environ.copy()

    if not os.path.exists(path):
        return env

    with open(path) as f:
        for raw in f:
            parsed = parse_env_line(raw)
            if not parsed:
                continue

            key, value = parsed

            # Merge known environment for expansion
            expanded = expand_vars(value, {**base, **env})
            env[key] = expanded

    return env