FROM confluentinc/cp-zookeeper:latest

# Set environment variables
ENV ZOOKEEPER_CLIENT_PORT=2181
ENV ZOOKEEPER_TICK_TIME=2000

# Add health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD echo 'ruok' | nc localhost 2181 | grep imok || exit 1

# Expose port
EXPOSE 2181