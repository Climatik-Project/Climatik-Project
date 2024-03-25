# Use a Python base image
FROM python:3.11-slim

# Set the working directory
WORKDIR /app

# Add requirements file
ADD requirements.txt .

# Install the required dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Add the operator code
ADD operator.py .

# Set the command to run the operator
CMD ["kopf", "run", "operator.py"]