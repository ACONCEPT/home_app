# Build stage
FROM instrumentisto/flutter:3 AS build-env

# Set up work directory
WORKDIR /app

# Copy the app files
COPY pubspec.* ./
RUN flutter pub get

COPY lib ./lib
COPY web ./web

# Build the app for web
RUN flutter build web --release

# Run stage - Use Nginx to serve the web app
FROM nginx:alpine

# Copy built web files
COPY --from=build-env /app/build/web /usr/share/nginx/html

# Copy custom nginx config
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
