FROM redis:7-alpine

# Set environment variables
ENV REDIS_PASSWORD=""
ENV REDIS_PORT=6379

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
  CMD redis-cli ping || exit 1

# Expose port
EXPOSE 6379