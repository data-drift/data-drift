from setuptools import setup, find_packages

setup(
    name="datagit",
    version="0.8",
    packages=find_packages(),
    author="Sammy Teillet",
    author_email="sammy.teillet@gmail.com",
    description="Git based metric store",
    long_description=open('README.md').read(),
    long_description_content_type='text/markdown',
    url="https://github.com/data-drift/datagit",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
    install_requires=[
        'pandas',
        'PyGithub'
    ],
)
