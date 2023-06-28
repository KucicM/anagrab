# Use the official Nginx base image
FROM nginx

# Copy the code files to the Nginx document root directory
COPY . /usr/share/nginx/html

# Create a volume for the Nginx document root directory
VOLUME /usr/share/nginx/html

