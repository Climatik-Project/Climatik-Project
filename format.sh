#!/bin/bash

# Check if YAPF is installed
if ! command -v yapf &> /dev/null
then
    echo "YAPF is not installed. Please install it by running 'pip install yapf' and try again."
    exit 1
fi

# Specify the style you want to use for YAPF. Examples: google, pep8, facebook, etc.
STYLE="pep8"

# Find all Python files in the current directory and its subdirectories,
# and run YAPF on them to format them in place with the specified style.
find python -type f -name "*.py"  -exec yapf --in-place --style=$STYLE "{}" +

echo "All Python files have been formatted with YAPF using the $STYLE style."

