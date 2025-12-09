"""
Commands module.
Contains Docker, migration, and lint command implementations.
"""

import os
import sys
import subprocess
from typing import Any


# ANSI Color Codes
RESET = "\033[0m"
BOLD = "\033[1m"
DIM = "\033[2m"
SYSTEM_COLOR = "\033[38;5;214m"
ERROR_COLOR = "\033[38;5;196m"


# Docker command functions
def docker_up_dev():
    """Start Docker development environment."""
    print("\n🐳 Starting Docker dev environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.dev.yml", "up", "-d"])


def docker_down_dev():
    """Stop Docker development environment."""
    print("\n🐳 Stopping Docker dev environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.dev.yml", "down"])


def docker_up_prod():
    """Start Docker production environment."""
    print("\n🚀 Starting Docker prod environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.prod.yml", "up", "-d"])


def docker_down_prod():
    """Stop Docker production environment."""
    print("\n🛑 Stopping Docker prod environment...")
    subprocess.run(
        [
            "docker",
            "compose",
            "-f",
            "docker-compose.prod.yml",
            "down",
            "--remove-orphans",
        ]
    )


def run_migrate(args):
    """Run migration command with parsed environment variables."""
    from scripts.env_parser import load_dotenv
    parsed_env = load_dotenv(".env")

    if args.subcommand == "create":
        # Special handling for migrate create command
        # Use fixed path as specified: /home/arkochan/Repositories/V1/db/pg/migrations
        migrations_dir = "/home/arkochan/Repositories/V1/db/pg/migrations"

        cwd = os.path.dirname(os.path.abspath(__file__))
        # Build the migrate create command with specific parameters
        absolute_dir = os.path.join(cwd, migrations_dir)
        cmd = [
            "migrate",
            "create",
            "-ext",
            "sql",
            "-dir",
            absolute_dir,
            "-seq",
            args.file_name,
        ]

        print(
            f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Running migration create: {' '.join(cmd)}{RESET}"
        )

        try:
            result = subprocess.run(cmd, check=True)
            print(f"{SYSTEM_COLOR}[system]{RESET} Migration created successfully")
        except subprocess.CalledProcessError as e:
            print(
                f"{ERROR_COLOR}[error]{RESET} Migration create failed with exit code {e.returncode}"
            )
            sys.exit(e.returncode)
        except FileNotFoundError:
            print(
                f"{ERROR_COLOR}[error]{RESET} migrate command not found. Make sure migrate is installed."
            )
            sys.exit(1)
    else:
        # Get migration source and database URL from environment for other commands
        migrations_source = parsed_env.get("MIGRATIONS")
        database_url = parsed_env.get("DATABASE_URL")

        if not migrations_source:
            print(
                f"{ERROR_COLOR}[error]{RESET} MIGRATIONS environment variable not found in .env"
            )
            sys.exit(1)
        if not database_url:
            print(
                f"{ERROR_COLOR}[error]{RESET} DATABASE_URL environment variable not found in .env"
            )
            sys.exit(1)

        cwd = os.path.dirname(os.path.abspath(__file__))
        # Build the migrate create command with specific parameters
        absolute_dir = os.path.join(cwd, migrations_source)
        # Build the migrate command
        cmd = ["migrate", "-path", absolute_dir, "-database", database_url]

        # Add the subcommand and any additional arguments
        if args.subcommand == "version":
            cmd.extend([args.subcommand] + getattr(args, "additional_args", []))
        else:
            cmd.extend([args.subcommand] + getattr(args, "additional_args", []))

        print(
            f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Running migration: {' '.join(cmd)}{RESET}"
        )

        try:
            result = subprocess.run(cmd, cwd=migrations_source, check=True)
            print(f"{SYSTEM_COLOR}[system]{RESET} Migration completed successfully")
        except subprocess.CalledProcessError as e:
            print(
                f"{ERROR_COLOR}[error]{RESET} Migration failed with exit code {e.returncode}"
            )
            sys.exit(e.returncode)
        except FileNotFoundError:
            print(
                f"{ERROR_COLOR}[error]{RESET} migrate command not found. Make sure migrate is installed."
            )
            sys.exit(1)


def run_lint(args):
    """Run linters for Go and Bun services."""
    if args.subcommand == "go":
        print(f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Running Go linter...{RESET}")
        try:
            result = subprocess.run(["make", "lint"], cwd="apps/go-app", check=True)
            print(f"{SYSTEM_COLOR}[system]{RESET} Go lint completed successfully")
        except subprocess.CalledProcessError as e:
            print(
                f"{ERROR_COLOR}[error]{RESET} Go lint failed with exit code {e.returncode}"
            )
            sys.exit(e.returncode)
        except FileNotFoundError:
            print(
                f"{ERROR_COLOR}[error]{RESET} make command not found or Makefile doesn't exist in apps/go-app"
            )
            sys.exit(1)
    elif args.subcommand == "bun":
        print(f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Running Bun linter...{RESET}")
        try:
            result = subprocess.run(["bun", "lint"], cwd="apps/next-app", check=True)
            print(f"{SYSTEM_COLOR}[system]{RESET} Bun lint completed successfully")
        except subprocess.CalledProcessError as e:
            print(
                f"{ERROR_COLOR}[error]{RESET} Bun lint failed with exit code {e.returncode}"
            )
            sys.exit(e.returncode)
        except FileNotFoundError:
            print(
                f"{ERROR_COLOR}[error]{RESET} bun command not found or lint script doesn't exist in apps/next-app"
            )
            sys.exit(1)
    else:
        print(f"{ERROR_COLOR}[error]{RESET} Please specify a linter: go or bun")
        sys.exit(1)


def run_init():
    """Initialize the repository with all necessary dependencies and tools."""
    import subprocess
    import os

    print(f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Initializing repository...{RESET}")

    # Change to the go-app directory
    go_app_dir = "apps/go-app"

    # Run go mod tidy
    print(f"{SYSTEM_COLOR}[system]{RESET} Running go mod tidy...")
    try:
        result = subprocess.run(["go", "mod", "tidy"], cwd=go_app_dir, check=True)
        print(f"{SYSTEM_COLOR}[system]{RESET} go mod tidy completed successfully")
    except subprocess.CalledProcessError as e:
        print(f"{ERROR_COLOR}[error]{RESET} go mod tidy failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"{ERROR_COLOR}[error]{RESET} go command not found. Please install Go.")
        sys.exit(1)

    # Install swag
    print(f"{SYSTEM_COLOR}[system]{RESET} Installing swag...")
    try:
        result = subprocess.run(["go", "install", "github.com/swaggo/swag/cmd/swag@latest"], check=True)
        print(f"{SYSTEM_COLOR}[system]{RESET} Swag installed successfully")
    except subprocess.CalledProcessError as e:
        print(f"{ERROR_COLOR}[error]{RESET} Swag installation failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"{ERROR_COLOR}[error]{RESET} go command not found. Please install Go.")
        sys.exit(1)

    # Run swag init
    print(f"{SYSTEM_COLOR}[system]{RESET} Running swag init...")
    try:
        result = subprocess.run(["swag", "init", "-g", "cmd/api/main.go"], cwd=go_app_dir, check=True)
        print(f"{SYSTEM_COLOR}[system]{RESET} Swag init completed successfully")
    except subprocess.CalledProcessError as e:
        print(f"{ERROR_COLOR}[error]{RESET} Swag init failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"{ERROR_COLOR}[error]{RESET} swag command not found. Make sure it's properly installed.")
        sys.exit(1)

    # Install air
    print(f"{SYSTEM_COLOR}[system]{RESET} Installing air...")
    try:
        result = subprocess.run(["go", "install", "github.com/air-verse/air@latest"], check=True)
        print(f"{SYSTEM_COLOR}[system]{RESET} Air installed successfully")
    except subprocess.CalledProcessError as e:
        print(f"{ERROR_COLOR}[error]{RESET} Air installation failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"{ERROR_COLOR}[error]{RESET} go command not found. Please install Go.")
        sys.exit(1)

    # Run bun install
    print(f"{SYSTEM_COLOR}[system]{RESET} Running bun install...")
    try:
        result = subprocess.run(["bun", "install"], cwd="apps/next-app", check=True)
        print(f"{SYSTEM_COLOR}[system]{RESET} Bun install completed successfully")
    except subprocess.CalledProcessError as e:
        print(f"{ERROR_COLOR}[error]{RESET} Bun install failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"{ERROR_COLOR}[error]{RESET} bun command not found. Please install Bun.")
        sys.exit(1)

    print(f"{SYSTEM_COLOR}[system]{RESET} {BOLD}Repository initialization completed successfully!{RESET}")