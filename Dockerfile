# Use the official Golang image as a base image
FROM golang:latest
WORKDIR /app
EXPOSE 8080
COPY google-chrome-stable_114.0.5735.90-1_amd64.deb .
RUN apt-get update && apt-get install -y \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libu2f-udev \
    libvulkan1 \
    libx11-6 \
    libxcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxkbcommon0 \
    libxrandr2 \
    xdg-utils
RUN dpkg -i google-chrome-stable_114.0.5735.90-1_amd64.deb && apt-get install -f
COPY go.mod .
RUN go mod download
COPY . .
CMD ["go", "run", "main/main.go"]
