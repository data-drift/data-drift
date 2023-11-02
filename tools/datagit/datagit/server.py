import http.server
import os
from threading import Timer
import webbrowser
import pkg_resources
import platform
import socket
import socketserver
import subprocess
import sys


def start_server(open_browser_url="/tables"):
    if is_port_in_use(9740) or is_port_in_use(9741):
        print("Server(s) already running on port 9740 or 9741. Exiting.")
        sys.exit()

    if platform.system() == "Darwin":
        if platform.machine().startswith("arm"):
            binary_path = pkg_resources.resource_filename(
                "datagit", "bin/datadrift-mac-m1"
            )
        else:
            binary_path = pkg_resources.resource_filename(
                "datagit", "bin/datadrift-mac-intel"
            )
    elif platform.system() == "Linux":
        binary_path = pkg_resources.resource_filename("datagit", "bin/datadrift-linux")
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
            print(f"Request path: {self.path}")
            # If the requested URL maps to an existing file, serve that.
            if os.path.exists(self.translate_path(self.path)):
                super().do_GET()
                return

            # Otherwise, serve the main index.html file.
            self.path = "index.html"
            super().do_GET()

    httpd = socketserver.TCPServer(("", PORT), Handler)

    try:
        print(f"Serving directory '{DIRECTORY}' on port {PORT}")
        url = f"http://localhost:{PORT}{open_browser_url}"
        print("Opening browser...", url)

        def open_url():
            webbrowser.open(url)

        Timer(1, open_url).start()
        httpd.serve_forever()
        server_process.wait()

    except KeyboardInterrupt:
        print("Shutting down servers...")
        httpd.shutdown()
        print("Httpd shut down")
        server_process.terminate()
        print("Server down")
        sys.exit()


def is_port_in_use(port):
    """Check if a given port is in use."""
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        return s.connect_ex(("localhost", port)) == 0
