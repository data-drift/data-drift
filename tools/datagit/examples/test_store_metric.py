import sys  # noqa
sys.path.insert(1, "../datagit")  # noqa
from datagit.github_connector import store_metric
import os
from dotenv import load_dotenv
from github import Github
import pandas as pd
import random

# Load environment variables from .env file
load_dotenv()
#

# Get GitHub token from environment variable
gh_token = os.getenv("GH_TOKEN")
if gh_token is None:
    print("GitHub token not found! Create a .env file a the root with a GH_TOKEN variable.")
    exit(1)

# Create GitHub client object
gh_client = Github(gh_token)

# Define test data
random_int = random.randint(0, 100)
data = {"unique_key": ["Alice", "Bob", "Charlie"], "age": [25, 30, random_int]}
df = pd.DataFrame(data)

# Define file path and assignees
file_path = os.getenv("PATH_FILE") or "gh_org/repo/path/to/file.csv"
branch = os.getenv("BRANCH") or "test"+str(random_int)
assignee = os.getenv("ASSIGNEE") or "Samox"

# Call store_metric function
store_metric(gh_client, df, file_path, assignees=[assignee], branch=branch)

print("Metric stored successfully!")
