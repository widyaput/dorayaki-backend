FROM golang:1.16.5-buster

# Making app folder, copy everything to it, and set working directory
RUN mkdir -p /app
COPY . /app
WORKDIR /app

# Update debian and install tools for easy debugging
RUN apt update

# Install and tidy golang dependency for server
RUN go mod tidy


# # Change script permission
# RUN chmod +x ./scripts/init-dev-server.sh

# Run the server
CMD ["bash", "./scripts/init-prod-server.sh"]

EXPOSE 8080
