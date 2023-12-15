from driftdb.version import version
from setuptools import find_packages, setup

setup(
    name="driftdb",
    version=version,
    packages=find_packages(),
    author="Sammy Teillet",
    author_email="sammy.teillet@gmail.com",
    description="Historical metric store",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    entry_points={"console_scripts": ["driftdb=driftdb.entrypoint:app"]},
    url="https://github.com/data-drift/data-drift/tree/main/tools/driftdb",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    package_data={
        "driftdb": [
            "bin/datadrift-mac-m1",
            "bin/datadrift-mac-intel",
            "bin/datadrift-linux",
            "bin/frontend/dist/**",
        ],
    },
    python_requires=">=3.6",
    install_requires=["pandas", "PyGithub", "click", "GitPython", "Faker", "typer"],
)
