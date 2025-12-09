import os
import sys
import signal
import time
import threading
import subprocess
from typing import Optional, Dict, List
import datetime


# ANSI Color Codes
RESET = "\033[0m"
BOLD = "\033[1m"
DIM = "\033[2m"
GO_COLOR = "\033[38;5;81m"
NEXT_COLOR = "\033[38;5;120m"
SYSTEM_COLOR = "\033[38;5;214m"
ERROR_COLOR = "\033[38;5;196m"
TIME_COLOR = "\033[38;5;244m"


class ServiceManager:
    """Manages a single service process with lifecycle operations."""

    def __init__(
        self,
        name: str,
        cmd: List[str],
        cwd: str,
        color: str,
        env: Optional[Dict[str, str]] = None,
    ):
        self.name = name
        self.cmd = cmd
        self.cwd = cwd
        self.color = color
        self.env = env or {}
        self.pid: Optional[int] = None
        self.status = "stopped"
        self.process: Optional[subprocess.Popen] = None
        self.lock = threading.Lock()

    def start(self) -> bool:
        """Start the service process."""
        with self.lock:
            if self.process and self.process.poll() is None:
                return True

            self.status = "starting"
            self._print_status(f"Starting {self.name}...")

            # Merge system env with parsed .env vars
            env = os.environ.copy()
            env.update(self.env)

            # First, do a quick check with a blocking run to catch startup errors
            try:
                startup_check = subprocess.run(
                    self.cmd,
                    cwd=self.cwd,
                    env=env,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    text=True,
                    timeout=2,  # Short timeout to catch immediate failures
                )
                if startup_check.returncode != 0:
                    # Command fails immediately
                    self.status = "stopped"
                    self._print_error(
                        f"{self.name} failed to start (exit code: {startup_check.returncode})"
                    )
                    print(f"{ERROR_COLOR}[{self.name}]{RESET} Startup output:")
                    for line in startup_check.stdout.strip().split("\n"):
                        if line.strip():
                            print(f"{ERROR_COLOR}[{self.name}]{RESET} {line}")
                    return False
            except subprocess.TimeoutExpired:
                # This means the command started successfully and is still running after 2s
                # This is expected for long-running services
                pass
            except FileNotFoundError as e:
                self.status = "stopped"
                self._print_error(f"Command not found: {self.cmd[0]} - {e}")
                return False
            except Exception as e:
                self.status = "stopped"
                self._print_error(f"Error during startup check: {e}")
                return False

            # Now start the actual long-running process
            try:
                self.process = subprocess.Popen(
                    self.cmd,
                    cwd=self.cwd,
                    env=env,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    text=True,
                    start_new_session=True,
                )
                self.pid = self.process.pid
                time.sleep(2)

                if self.process.poll() is None:
                    self.status = "running"
                    self._print_status(f"{self.name} started (PID: {self.pid})")

                    # Start log forwarding thread
                    thread = threading.Thread(target=self._log_forwarder, daemon=True)
                    thread.start()
                    return True
                else:
                    self.status = "stopped"
                    exit_code = self.process.returncode
                    self._print_error(
                        f"{self.name} failed to start (exit code: {exit_code})"
                    )
                    return False

            except Exception as e:
                self.status = "stopped"
                self._print_error(f"Error starting {self.name}: {e}")
                return False

    def stop(self) -> None:
        """Stop service and its entire process tree."""
        with self.lock:
            if not self.process or self.process.poll() is not None:
                self.pid = None
                self.process = None
                self.status = "stopped"
                return

            self._print_status(f"Stopping {self.name} (PID: {self.pid})...")

            try:
                # Kill entire process group using negative PID
                if self.pid is None:
                    self._print_error(f"No PID found for {self.name}, cannot stop")
                    self.process = None
                    self.pid = None
                    self.status = "stopped"
                    return
                # On macOS, we need to use os.killpg with the process group ID
                pgid = os.getpgid(self.pid)
                os.killpg(pgid, signal.SIGTERM)

                # Wait for graceful shutdown (5s max)
                for _ in range(50):
                    if self.process.poll() is not None:
                        break
                    time.sleep(0.1)
                else:
                    # Force kill if still alive
                    self._print_status(f"Force killing {self.name}...")
                    os.killpg(pgid, signal.SIGKILL)
                    time.sleep(0.5)

            except ProcessLookupError:
                pass
            except OSError as e:
                # Handle case where process group doesn't exist
                if e.errno == 3:  # No such process
                    pass
                else:
                    self._print_error(f"Error stopping {self.name}: {e}")
            except Exception as e:
                self._print_error(f"Error stopping {self.name}: {e}")

            self.process = None
            self.pid = None
            self.status = "stopped"

    def restart(self) -> bool:
        """Restart the service."""
        self._print_status(f"Restarting {self.name}...")
        self.status = "restarting"
        self.stop()
        time.sleep(0.5)
        return self.start()

    def _log_forwarder(self) -> None:
        """Forward and format logs from service stdout."""
        if not self.process or not self.process.stdout:
            return

        for line in iter(self.process.stdout.readline, ""):
            if line and not self.process.poll():
                timestamp = datetime.datetime.now().strftime("%H:%M:%S")
                gutter = "[go]  " if self.name == "Go" else "[next]"
                print(
                    f"{TIME_COLOR}{timestamp}{RESET} {self.color}{gutter}{RESET} {line.rstrip()}"
                )
                sys.stdout.flush()

    def _print_status(self, msg: str) -> None:
        """Print system status message."""
        print(f"{SYSTEM_COLOR}[system]{RESET} {BOLD}{msg}{RESET}")
        sys.stdout.flush()

    def _print_error(self, msg: str) -> None:
        """Print error message."""
        print(f"{ERROR_COLOR}[system]{RESET} {msg}")
        sys.stdout.flush()

