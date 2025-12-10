import os
import sys
import signal
import time
import threading
import subprocess
from typing import Optional, Dict, List

from scripts.loggers import log_system, log_error, log_service_output, GO_COLOR, NEXT_COLOR


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
            log_system(f"Starting {self.name}...")

            # Merge system env with parsed .env vars
            env = os.environ.copy()
            env.update(self.env)

            # Start the actual long-running process directly (no initial check to avoid double execution)
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

                # Start log forwarding thread immediately to capture all output
                thread = threading.Thread(target=self._log_forwarder, daemon=True)
                thread.start()

                # Wait a moment to see if process stays alive
                time.sleep(2)

                if self.process.poll() is None:
                    self.status = "running"
                    log_system(f"{self.name} started (PID: {self.pid})")
                    return True
                else:
                    # Wait briefly to let log forwarder process any remaining output
                    time.sleep(0.5)
                    self.status = "stopped"
                    exit_code = self.process.returncode
                    log_error(
                        f"{self.name} failed to start (exit code: {exit_code})"
                    )
                    return False

            except FileNotFoundError as e:
                self.status = "stopped"
                log_error(f"Command not found: {self.cmd[0]} - {e}")
                return False
            except Exception as e:
                self.status = "stopped"
                log_error(f"Error starting {self.name}: {e}")
                return False

    def stop(self) -> None:
        """Stop service and its entire process tree."""
        with self.lock:
            if not self.process or self.process.poll() is not None:
                self.pid = None
                self.process = None
                self.status = "stopped"
                return

            log_system(f"Stopping {self.name} (PID: {self.pid})...")

            try:
                # Kill entire process group using negative PID
                if self.pid is None:
                    log_error(f"No PID found for {self.name}, cannot stop")
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
                    log_system(f"Force killing {self.name}...")
                    os.killpg(pgid, signal.SIGKILL)
                    time.sleep(0.5)

            except ProcessLookupError:
                pass
            except OSError as e:
                # Handle case where process group doesn't exist
                if e.errno == 3:  # No such process
                    pass
                else:
                    log_error(f"Error stopping {self.name}: {e}")
            except Exception as e:
                log_error(f"Error stopping {self.name}: {e}")

            self.process = None
            self.pid = None
            self.status = "stopped"

    def restart(self) -> bool:
        """Restart the service."""
        log_system(f"Restarting {self.name}...")
        self.status = "restarting"
        self.stop()
        time.sleep(0.5)
        return self.start()

    def _log_forwarder(self) -> None:
        """Forward and format logs from service stdout."""
        if not self.process or not self.process.stdout:
            return

        for line in iter(self.process.stdout.readline, ""):
            if line:  # Remove the poll() check to capture logs even after process exits
                color = GO_COLOR if self.name.lower() == "go" else NEXT_COLOR
                log_service_output(self.name, line.rstrip(), color)
                sys.stdout.flush()

        # After exiting the loop, check if process exited with error and print any remaining info
        if self.process and self.process.poll() is not None and self.process.returncode != 0:
            log_error(
                f"{self.name} exited with code {self.process.returncode}"
            )


