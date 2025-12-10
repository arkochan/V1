"""
Terminal User Interface module.
Provides a terminal UI for managing services with keyboard shortcuts.
"""

import sys
import select
import termios
import tty
import signal
import time

from scripts.loggers import (
    log_system,
    SYSTEM_COLOR,
    ERROR_COLOR,
    BOLD,
    RESET,
    DIM,
    GO_COLOR,
    NEXT_COLOR,
)


class TUI:
    """Terminal User Interface for service management."""

    def __init__(self, go_enabled=True, next_enabled=True, service_manager_class=None):
        if service_manager_class is None:
            from scripts.service_manager import ServiceManager

            self.service_manager_class = ServiceManager
        else:
            self.service_manager_class = service_manager_class

        # Parse .env file with variable expansion
        from scripts.env_parser import load_dotenv

        self.parsed_env = load_dotenv(".env")

        if self.parsed_env:
            log_system(f"{DIM}Loaded {len(self.parsed_env)} variables from .env{RESET}")
        else:
            log_system(f"{DIM}No .env file found or empty{RESET}")

        # Create services based on arguments
        self.services = {}
        self.go_enabled = go_enabled
        self.next_enabled = next_enabled

        if go_enabled:
            self.services["go"] = self.service_manager_class(
                name="Go",
                cmd=["make", "dev"],
                cwd="apps/go-app",
                color=GO_COLOR,
                env=self.parsed_env,  # All parsed env vars available to Go
            )

        if next_enabled:
            self.services["next"] = self.service_manager_class(
                name="Next.js",
                cmd=["bun", "run", "dev"],
                cwd="apps/next-app",
                color=NEXT_COLOR,
                env=self.parsed_env,  # All parsed env vars available to Next.js
            )

        self.last_key = ""
        self.last_key_time = 0.0
        self.shutting_down = False

    def setup_terminal(self) -> None:
        """Configure terminal for immediate input."""
        self.fd = sys.stdin.fileno()
        self.old_settings = termios.tcgetattr(self.fd)
        tty.setcbreak(self.fd)
        sys.stdout.write("\033[?25l")  # Hide cursor
        sys.stdout.flush()

    def restore_terminal(self) -> None:
        """Restore terminal to original state."""
        sys.stdout.write("\033[?25h")  # Show cursor
        sys.stdout.flush()
        termios.tcsetattr(self.fd, termios.TCSADRAIN, self.old_settings)

    def print_header(self) -> None:
        """Display current service status header."""
        print(
            f"{BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{RESET}"
        )

        # Print Go status if enabled
        if self.go_enabled:
            print(f"{BOLD}  Go Backend: {RESET}", end="")
            symbols = {
                "running": f"{GO_COLOR}●{RESET} running  ",
                "starting": f"{SYSTEM_COLOR}◐{RESET} starting ",
                "stopped": f"{ERROR_COLOR}○{RESET} stopped  ",
                "restarting": f"{SYSTEM_COLOR}↻{RESET} restarting ",
            }
            go_status = self.services["go"].status
            print(symbols.get(go_status, "unknown "), end="")

        # Print Next status if enabled
        if self.next_enabled:
            print(f"{BOLD}Next.js: {RESET}", end="")
            symbols = {
                "running": f"{NEXT_COLOR}●{RESET} running",
                "starting": f"{SYSTEM_COLOR}◐{RESET} starting",
                "stopped": f"{ERROR_COLOR}○{RESET} stopped",
                "restarting": f"{SYSTEM_COLOR}↻{RESET} restarting",
            }
            next_status = self.services["next"].status
            print(symbols.get(next_status, "unknown"))

        # Print controls based on which services are enabled
        controls = []
        if self.go_enabled and self.next_enabled:
            controls.append("g→r: restart Go")
            controls.append("n→r: restart Next")
            controls.append("r: restart all")
        elif self.go_enabled:
            controls.append("g→r: restart Go")
        elif self.next_enabled:
            controls.append("n→r: restart Next")

        controls.append("Ctrl+C: quit")
        print(f"{DIM}  {'  |  '.join(controls)}{RESET}")
        print(
            f"{BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{RESET}"
        )
        sys.stdout.flush()

    def handle_input(self) -> None:
        """Process keyboard commands."""
        if select.select([sys.stdin], [], [], 0.1)[0]:
            char = sys.stdin.read(1)
            current_time = time.time()

            # Reset chord on timeout
            if self.last_key and (current_time - self.last_key_time > 1.0):
                self.last_key = ""

            chord = self.last_key + char

            if chord == "gr" and self.go_enabled:
                self.last_key = ""
                print("\033[2J\033[H")  # Clear screen
                self.print_header()
                self.services["go"].restart()
            elif chord == "nr" and self.next_enabled:
                self.last_key = ""
                print("\033[2J\033[H")
                self.print_header()
                self.services["next"].restart()
            elif chord == "r" and not self.last_key:
                self.last_key = ""
                print("\033[2J\033[H")
                self.print_header()
                self.restart_all()
            elif char in ("g", "n"):
                # Only accept if the corresponding service is enabled
                if (char == "g" and self.go_enabled) or (
                    char == "n" and self.next_enabled
                ):
                    self.last_key = char
                    self.last_key_time = current_time
            else:
                self.last_key = ""

    def restart_all(self) -> None:
        """Restart both services."""
        log_system(f"{BOLD}Restarting all services...{RESET}")
        for service in self.services.values():
            service.restart()
        log_system(f"{DIM}All services restarted{RESET}")
        sys.stdout.flush()

    def cleanup(self) -> None:
        """Graceful shutdown of services and terminal restoration."""
        if self.shutting_down:
            return
        self.shutting_down = True

        print()
        log_system(f"{BOLD}Shutting down...{RESET}")

        for service in self.services.values():
            service.stop()

        self.restore_terminal()
        time.sleep(0.1)  # Ensure terminal is fully restored

        log_system(f"{DIM}All services stopped. Goodbye!{RESET}")
        sys.stdout.flush()

    def run(self) -> None:
        """Main TUI execution loop."""

        def signal_handler(sig, frame):
            self.cleanup()
            sys.exit(0)

        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)

        self.setup_terminal()

        print("\033[2J\033[H")  # Clear screen
        self.print_header()

        log_system(f"{BOLD}Starting services...{RESET}")
        for service in self.services.values():
            service.start()

        try:
            while not self.shutting_down:
                self.handle_input()

                # Health monitoring
                for name, service in self.services.items():
                    if service.status == "running" and service.process:
                        if service.process.poll() is not None:
                            # Process has stopped unexpectedly, show recent logs
                            if service.recent_logs:
                                print(
                                    f"{ERROR_COLOR}[{service.name}]{RESET} Recent logs before exit:"
                                )
                                for log_line in service.recent_logs:
                                    print(
                                        f"{ERROR_COLOR}[{service.name}]{RESET} {log_line}"
                                    )

                            service.status = "stopped"
                            exit_code = service.process.returncode
                            from scripts.loggers import log_simple_error

                            log_simple_error(
                                f"{service.name} stopped unexpectedly (exit code: {exit_code})"
                            )

                time.sleep(0.1)
        finally:
            # Ensure cleanup runs even on unexpected errors
            self.cleanup()
