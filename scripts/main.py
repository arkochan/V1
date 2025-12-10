#!/usr/bin/env python3
"""
Development environment manager for Go backend and Next.js frontend.
Provides a terminal UI for managing services with keyboard shortcuts.
No external dependencies - uses only Python standard library.
"""

import os
import sys
import signal
import time
import argparse
from typing import Optional


def main():
    """Main command dispatcher."""
    parser = argparse.ArgumentParser(description="Development environment manager")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # TUI subcommands
    tui_parser = subparsers.add_parser("tui", help="Start TUI development environment")
    tui_group = tui_parser.add_mutually_exclusive_group()
    tui_group.add_argument("--go", action="store_true", help="Only start Go service")
    tui_group.add_argument(
        "--next", action="store_true", help="Only start Next.js service"
    )
    tui_group.add_argument(
        "--dev", action="store_true", help="Start both services (default)"
    )

    # Dev command (alias for TUI)
    dev_parser = subparsers.add_parser("dev", help="Start development environment")
    dev_group = dev_parser.add_mutually_exclusive_group()
    dev_group.add_argument("--go", action="store_true", help="Only start Go service")
    dev_group.add_argument(
        "--next", action="store_true", help="Only start Next.js service"
    )
    dev_group.add_argument(
        "--dev", action="store_true", help="Start both services (default)"
    )

    # Simple go command to start only Go service
    subparsers.add_parser("go", help="Start only Go service")

    # Simple next command to start only Next.js service
    subparsers.add_parser("next", help="Start only Next.js service")

    # Docker commands
    subparsers.add_parser("up", help="Start Docker dev environment")
    subparsers.add_parser("down", help="Stop Docker dev environment")
    subparsers.add_parser("prod-up", help="Start Docker prod environment")
    subparsers.add_parser("prod-down", help="Stop Docker prod environment")

    # Migration command
    from .commands import run_migrate
    migrate_parser = subparsers.add_parser("migrate", help="Run database migrations")
    migrate_subparsers = migrate_parser.add_subparsers(
        dest="subcommand", help="Migration subcommands"
    )

    # Add subcommands for migrate
    migrate_up = migrate_subparsers.add_parser("up", help="Apply migrations")
    migrate_up.add_argument(
        "additional_args", nargs="*", help="Additional arguments for up command"
    )

    migrate_down = migrate_subparsers.add_parser("down", help="Rollback migrations")
    migrate_down.add_argument(
        "additional_args", nargs="*", help="Additional arguments for down command"
    )

    migrate_create = migrate_subparsers.add_parser(
        "create", help="Create new migration"
    )
    migrate_create.add_argument("file_name", help="Name for the new migration file")

    migrate_drop = migrate_subparsers.add_parser("drop", help="Drop all migrations")
    migrate_drop.add_argument(
        "additional_args", nargs="*", help="Additional arguments for drop command"
    )

    migrate_force = migrate_subparsers.add_parser("force", help="Force version")
    migrate_force.add_argument(
        "additional_args", nargs="*", help="Additional arguments for force command"
    )

    migrate_version = migrate_subparsers.add_parser(
        "version", help="Get migration version"
    )
    migrate_version.add_argument(
        "additional_args", nargs="*", help="Additional arguments for version command"
    )

    # Lint command
    lint_parser = subparsers.add_parser("lint", help="Run linters")
    lint_subparsers = lint_parser.add_subparsers(
        dest="subcommand", help="Lint subcommands"
    )

    # Add subcommands for lint
    lint_go = lint_subparsers.add_parser("go", help="Run Go linter")
    lint_bun = lint_subparsers.add_parser("bun", help="Run Bun linter")

    # Init command
    subparsers.add_parser("init", help="Initialize the repository with dependencies and tools")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    try:
        if args.command in ("dev", "tui"):
            # Determine which services to enable
            from scripts.tui import TUI

            go_enabled = getattr(args, "go", False) or not (
                getattr(args, "next", False)
            )
            next_enabled = getattr(args, "next", False) or not (
                getattr(args, "go", False)
            )

            tui = TUI(go_enabled=go_enabled, next_enabled=next_enabled)
            tui.run()
        elif args.command == "go":
            # Start only the Go service without TUI
            from scripts.env_parser import load_dotenv
            from scripts.service_manager import ServiceManager
            from scripts.loggers import log_system, log_simple_error, BOLD, DIM

            # ANSI Color Codes
            RESET = "\033[0m"
            GO_COLOR = "\033[38;5;81m"

            parsed_env = load_dotenv(".env")
            go_service = ServiceManager(
                name="Go",
                cmd=["make", "dev"],
                cwd="apps/go-app",
                color=GO_COLOR,
                env=parsed_env,
            )
            log_system(f"{BOLD}Starting Go service...{RESET}")
            if go_service.start():
                log_system(f"{BOLD}Go service started. Press Ctrl+C to stop.{RESET}")
                try:
                    while True:
                        time.sleep(1)
                        if go_service.process and go_service.process.poll() is not None:
                            log_simple_error("Go service stopped unexpectedly")
                            break
                except KeyboardInterrupt:
                    print(f"\n", end="")
                    log_system(f"{BOLD}Stopping Go service...{RESET}")
                    go_service.stop()
                    log_system(f"{DIM}Go service stopped.{RESET}")
        elif args.command == "next":
            # Start only the Next.js service without TUI
            from scripts.env_parser import load_dotenv
            from scripts.service_manager import ServiceManager
            from scripts.loggers import log_system, log_simple_error, BOLD, DIM

            # ANSI Color Codes
            RESET = "\033[0m"
            NEXT_COLOR = "\033[38;5;120m"

            parsed_env = load_dotenv(".env")
            next_service = ServiceManager(
                name="Next.js",
                cmd=["bun", "run", "dev"],
                cwd="apps/next-app",
                color=NEXT_COLOR,
                env=parsed_env,
            )
            log_system(f"{BOLD}Starting Next.js service...{RESET}")
            if next_service.start():
                log_system(f"{BOLD}Next.js service started. Press Ctrl+C to stop.{RESET}")
                try:
                    while True:
                        time.sleep(1)
                        if (
                            next_service.process
                            and next_service.process.poll() is not None
                        ):
                            log_simple_error("Next.js service stopped unexpectedly")
                            break
                except KeyboardInterrupt:
                    print(f"\n", end="")
                    log_system(f"{BOLD}Stopping Next.js service...{RESET}")
                    next_service.stop()
                    log_system(f"{DIM}Next.js service stopped.{RESET}")
        elif args.command == "up":
            from scripts.commands import docker_up_dev
            docker_up_dev()
        elif args.command == "down":
            from scripts.commands import docker_down_dev
            docker_down_dev()
        elif args.command == "prod-up":
            from scripts.commands import docker_up_prod
            docker_up_prod()
        elif args.command == "prod-down":
            from scripts.commands import docker_down_prod
            docker_down_prod()
        elif args.command == "migrate":
            if not args.subcommand:
                migrate_parser.print_help()
                sys.exit(1)
            from scripts.commands import run_migrate
            run_migrate(args)
        elif args.command == "lint":
            if not args.subcommand:
                print("Usage: ./run lint [go|bun]")
                sys.exit(1)
            from scripts.commands import run_lint
            run_lint(args)
        elif args.command == "init":
            from scripts.commands import run_init
            run_init()
        else:
            parser.print_help()
            sys.exit(1)
    except KeyboardInterrupt:
        print("\nInterrupted by user")
        sys.exit(0)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()