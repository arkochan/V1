"""
Commands module.
Contains Docker, migration, and lint command implementations.
"""

import os
import sys
import subprocess
from typing import Any

from scripts.loggers import log_system, log_simple_error, SYSTEM_COLOR, ERROR_COLOR, BOLD, RESET


# Docker command functions
def docker_up_dev():
    """Start Docker development environment."""
    log_system("🐳 Starting Docker dev environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.dev.yml", "up", "-d"])


def docker_down_dev():
    """Stop Docker development environment."""
    log_system("🐳 Stopping Docker dev environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.dev.yml", "down"])


def docker_up_prod():
    """Start Docker production environment."""
    log_system("🚀 Starting Docker prod environment...")
    subprocess.run(["docker", "compose", "-f", "docker-compose.prod.yml", "up", "-d"])


def docker_down_prod():
    """Stop Docker production environment."""
    log_system("🛑 Stopping Docker prod environment...")
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

        log_system(f"{BOLD}Running migration create: {' '.join(cmd)}{RESET}")

        try:
            result = subprocess.run(cmd, check=True)
            log_system("Migration created successfully")
        except subprocess.CalledProcessError as e:
            log_simple_error(f"Migration create failed with exit code {e.returncode}")
            sys.exit(e.returncode)
        except FileNotFoundError:
            log_simple_error("migrate command not found. Make sure migrate is installed.")
            sys.exit(1)
    else:
        # Get migration source and database URL from environment for other commands
        migrations_source = parsed_env.get("MIGRATIONS")
        database_url = parsed_env.get("DATABASE_URL")

        if not migrations_source:
            log_simple_error("MIGRATIONS environment variable not found in .env")
            sys.exit(1)
        if not database_url:
            log_simple_error("DATABASE_URL environment variable not found in .env")
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

        log_system(f"{BOLD}Running migration: {' '.join(cmd)}{RESET}")

        try:
            result = subprocess.run(cmd, cwd=migrations_source, check=True)
            log_system("Migration completed successfully")
        except subprocess.CalledProcessError as e:
            log_simple_error(f"Migration failed with exit code {e.returncode}")
            sys.exit(e.returncode)
        except FileNotFoundError:
            log_simple_error("migrate command not found. Make sure migrate is installed.")
            sys.exit(1)


def run_lint(args):
    """Run linters for Go and Bun services."""
    if args.subcommand == "go":
        log_system(f"{BOLD}Running Go linter...{RESET}")
        try:
            result = subprocess.run(["make", "lint"], cwd="apps/go-app", check=True)
            log_system("Go lint completed successfully")
        except subprocess.CalledProcessError as e:
            log_simple_error(f"Go lint failed with exit code {e.returncode}")
            sys.exit(e.returncode)
        except FileNotFoundError:
            log_simple_error("make command not found or Makefile doesn't exist in apps/go-app")
            sys.exit(1)
    elif args.subcommand == "bun":
        log_system(f"{BOLD}Running Bun linter...{RESET}")
        try:
            result = subprocess.run(["bun", "lint"], cwd="apps/next-app", check=True)
            log_system("Bun lint completed successfully")
        except subprocess.CalledProcessError as e:
            log_simple_error(f"Bun lint failed with exit code {e.returncode}")
            sys.exit(e.returncode)
        except FileNotFoundError:
            log_simple_error("bun command not found or lint script doesn't exist in apps/next-app")
            sys.exit(1)
    else:
        log_simple_error("Please specify a linter: go or bun")
        sys.exit(1)


def run_init():
    """Initialize the repository with all necessary dependencies and tools."""
    import subprocess
    import os

    log_system(f"{BOLD}Initializing repository...{RESET}")

    # Change to the go-app directory
    go_app_dir = "apps/go-app"

    # Run go mod tidy
    log_system("Running go mod tidy...")
    try:
        result = subprocess.run(["go", "mod", "tidy"], cwd=go_app_dir, check=True)
        log_system("go mod tidy completed successfully")
    except subprocess.CalledProcessError as e:
        log_simple_error(f"go mod tidy failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        log_simple_error("go command not found. Please install Go.")
        sys.exit(1)

    # Install swag
    log_system("Installing swag...")
    try:
        result = subprocess.run(["go", "install", "github.com/swaggo/swag/cmd/swag@latest"], check=True)
        log_system("Swag installed successfully")
    except subprocess.CalledProcessError as e:
        log_simple_error(f"Swag installation failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        log_simple_error("go command not found. Please install Go.")
        sys.exit(1)

    # Run swag init
    log_system("Running swag init...")
    try:
        result = subprocess.run(["swag", "init", "-g", "cmd/api/main.go"], cwd=go_app_dir, check=True)
        log_system("Swag init completed successfully")
    except subprocess.CalledProcessError as e:
        log_simple_error(f"Swag init failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        log_simple_error("swag command not found. Make sure it's properly installed.")
        sys.exit(1)

    # Install air
    log_system("Installing air...")
    try:
        result = subprocess.run(["go", "install", "github.com/air-verse/air@latest"], check=True)
        log_system("Air installed successfully")
    except subprocess.CalledProcessError as e:
        log_simple_error(f"Air installation failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        log_simple_error("go command not found. Please install Go.")
        sys.exit(1)

    # Run bun install
    log_system("Running bun install...")
    try:
        result = subprocess.run(["bun", "install"], cwd="apps/next-app", check=True)
        log_system("Bun install completed successfully")
    except subprocess.CalledProcessError as e:
        log_simple_error(f"Bun install failed with exit code {e.returncode}")
        sys.exit(e.returncode)
    except FileNotFoundError:
        log_simple_error("bun command not found. Please install Bun.")
        sys.exit(1)

    log_system(f"{BOLD}Repository initialization completed successfully!{RESET}")