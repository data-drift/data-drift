import http.server
import os
import platform
import socket
import socketserver
import subprocess
import sys
import webbrowser
from threading import Timer

import pkg_resources

from .logger import get_logger

logger = get_logger(__name__)


def start_server(open_browser_url="/tables"):
    if is_port_in_use(9740) or is_port_in_use(9741):
        logger.warn("Server(s) already running on port 9740 or 9741. Exiting.")
        sys.exit()

    if platform.system() == "Darwin":
        if platform.machine().startswith("arm"):
            binary_path = pkg_resources.resource_filename("driftdb", "bin/datadrift-mac-m1")
        else:
            binary_path = pkg_resources.resource_filename("driftdb", "bin/datadrift-mac-intel")
    elif platform.system() == "Linux":
        binary_path = pkg_resources.resource_filename("driftdb", "bin/datadrift-linux")
    else:
        # TODO: Update this path for other platforms (Linux, Windows, etc.)
        raise ValueError("Unsupported platform")

    # Get a copy of the current environment variables
    env = os.environ.copy()

    # Set the PORT environment variable
    env["PORT"] = "9740"

    server_process = subprocess.Popen(
        [binary_path],
        env=env,
    )

    PORT = 9741

    SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
    DIRECTORY = os.path.join(SCRIPT_DIR, "bin/frontend/dist")

    class Handler(http.server.SimpleHTTPRequestHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, directory=DIRECTORY, **kwargs)

        def do_GET(self):
            logger.info(f"Request path: {self.path}")
            # If the requested URL maps to an existing file, serve that.
            if os.path.exists(self.translate_path(self.path)):
                super().do_GET()
                return

            # Otherwise, serve the main index.html file.
            self.path = "index.html"
            super().do_GET()

    httpd = socketserver.TCPServer(("", PORT), Handler)

    try:
        logger.info(f"Serving directory '{DIRECTORY}' on port {PORT}")
        url = f"http://localhost:{PORT}{open_browser_url}"
        logger.info(f"Opening browser... {url}")

        def open_url():
            webbrowser.open(url)

        Timer(1, open_url).start()
        httpd.serve_forever()
        server_process.wait()

    except KeyboardInterrupt:
        logger.info("Shutting down servers...")
        httpd.shutdown()
        logger.info("Httpd shut down")
        server_process.terminate()
        logger.info("Server down")
        sys.exit()


def is_port_in_use(port):
    """Check if a given port is in use."""
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        return s.connect_ex(("localhost", port)) == 0
