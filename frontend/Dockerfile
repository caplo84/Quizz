# Use ARG to allow dynamic Node version (defaults to 20)
ARG NODE_VERSION=20
FROM node:${NODE_VERSION}-alpine

# Set the working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy app files and fix permissions
COPY . .

# 🛠️ Fix ownership of files
RUN chown -R node:node /app

# Switch to non-root user (recommended)
USER node

# Build the app
RUN npm run build

# Expose the Vite preview port
EXPOSE 5173

# Start the preview server
CMD ["npm", "run", "preview"]
