FROM postgres:15-alpine

# Set environment variables
ENV POSTGRES_DB=demo_app
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password

# Create init script directory
RUN mkdir -p /docker-entrypoint-initdb.d

# Copy initialization scripts
COPY ./migrations/001_initial_schema.up.sql /docker-entrypoint-initdb.d/01-schema.sql

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
  CMD pg_isready -U postgres -d demo_app || exit 1

# Expose port
EXPOSE 5432