# Use the base image with Node.js
FROM nikolaik/python-nodejs:latest as build

# Install Go
RUN wget https://dl.google.com/go/go1.18.3.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz \
    && rm go1.18.3.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Set up working directory for Go application
WORKDIR /app

# Install Air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Add Go bin to PATH
ENV PATH="/root/go/bin:${PATH}"

# Copy Go files and other necessary files into the container
COPY . .

# Install Node.js dependencies (including Tailwind CSS)
RUN npm install
RUN npm install -g postcss-cli

# Set the command to run Air
CMD ["sh", "-c", "ls -la /app && air"]