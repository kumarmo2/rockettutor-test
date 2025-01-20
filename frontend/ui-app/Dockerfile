# FROM node:20-alpine AS build-env
# WORKDIR /app
# COPY . .
# RUN npm run build
#
# FROM node:20-alpine
# WORKDIR /app
# COPY package.json package-lock.json ./
# COPY . .
# CMD ["npm", "run", "start"]


# Use Node.js for building the app
FROM node:20-alpine AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the application
RUN npm run build

# Use a lightweight nginx image to serve the app
FROM nginx:stable-alpine

# Copy the built app from the previous stage to nginx's web root
COPY --from=build /app/dist /usr/share/nginx/html

# Copy a custom nginx configuration file if needed (optional)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose the port nginx will run on
EXPOSE 5000

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

