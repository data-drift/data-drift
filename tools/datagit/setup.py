from setuptools import setup, find_packages

setup(
    name="driftdb",
    version="0.0.1-alpha5",
    packages=find_packages(),
    author="Sammy Teillet",
    author_email="sammy.teillet@gmail.com",
    description="Git based metric store",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    entry_points={"console_scripts": ["driftdb=datagit.cli:cli_entrypoint"]},
    url="https://github.com/data-drift/data-drift/tree/main/tools/datagit",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    package_data={
        "datagit": [
            "bin/data-drift-mac-m1",
            "bin/data-drift-mac-intel",
            "bin/frontend/dist/**",
        ],
    },
    python_requires=">=3.6",
    install_requires=["pandas", "PyGithub", "click", "GitPython", "Faker"],
)
