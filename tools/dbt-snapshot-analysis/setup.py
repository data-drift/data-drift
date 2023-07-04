from setuptools import setup, find_packages

setup(
    name="dbt_snapshot_analysis",
    version="0.2.11",
    packages=find_packages(),
    py_modules=["dbt_snapshot_analysis"],
    install_requires=["pandas", "plotly", "streamlit"],
    long_description=open("README.md", "r").read(),
    long_description_content_type="text/markdown",
    entry_points={
        "console_scripts": [
            "dbt_snapshot_analysis=dbt_snapshot_analysis.streamlit_entry_point:main"
        ]
    },
    author="Sammy Teillet",
    author_email="sammy.teillet@gmail.com",
    description="A package for analyzing snapshots",
    url="https://github.com/data-drift/data-drift/tree/main/tools/dbt-snapshot-analysis",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
)
