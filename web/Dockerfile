FROM node:22-alpine

WORKDIR /app

# Install dependencies first (caching)
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code
COPY . .

# Expose the dev server port
EXPOSE 10104

# Run the dev server
CMD ["npm", "run", "dev"]
