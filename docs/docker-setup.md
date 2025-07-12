# Docker Setup Guide

This guide explains how to set up and run the Great Nigeria Library Foundation using Docker.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Git
- At least 4GB RAM available for Docker
- At least 10GB free disk space

## Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/yerenwgventures/GreatNigeriaLibrary-Foundation.git
cd GreatNigeriaLibrary-Foundation
cp .env.example .env
```

### 2. Configure Environment

Edit the `.env` file with your settings:

```bash
# Required: Change these for security
DB_PASSWORD=your_secure_database_password
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters
REDIS_PASSWORD=your_secure_redis_password

# Optional: Customize ports if needed
SERVER_EXTERNAL_PORT=8080
DB_EXTERNAL_PORT=5433
REDIS_EXTERNAL_PORT=6380
```

### 3. Start the Application

```bash
# Production-like environment
docker-compose up -d

# Development environment with hot reloading
./scripts/docker-dev.sh start

# Development with admin tools
./scripts/docker-dev.sh start-tools
```

## Docker Compose Files

### Main Configuration (`docker-compose.yml`)

The main Docker Compose file includes:

- **PostgreSQL 15**: Primary database with health checks
- **Redis 7**: Caching and session storage
- **Application**: Go backend with multi-stage build
- **Nginx**: Reverse proxy (production profile)

### Development Override (`docker-compose.dev.yml`)

Additional services for development:

- **Hot Reloading**: Air for automatic rebuilds
- **Database Admin**: Adminer for database management
- **Redis Admin**: Redis Commander for cache management
- **Mail Catcher**: MailHog for email testing
- **Debug Ports**: pprof and other debugging tools

### Production Configuration (`docker-compose.prod.yml`)

Production optimizations:

- **Resource Limits**: CPU and memory constraints
- **Security**: No external port exposure for databases
- **Monitoring**: Prometheus and Grafana (optional)
- **High Availability**: Multiple app replicas
- **Performance Tuning**: Optimized database and Redis settings

## Service Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Nginx       │    │   Application   │    │   PostgreSQL    │
│  (Port 80/443)  │────│   (Port 8080)   │────│   (Port 5432)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                │
                       ┌─────────────────┐
                       │      Redis      │
                       │   (Port 6379)   │
                       └─────────────────┘
```

## Development Workflow

### Starting Development Environment

```bash
# Start basic development environment
./scripts/docker-dev.sh start

# Start with development tools
./scripts/docker-dev.sh start-tools
```

### Available Services (Development)

- **Application**: http://localhost:8081
- **Database Admin**: http://localhost:8082
- **Redis Admin**: http://localhost:8083
- **Mail Catcher**: http://localhost:8025

### Development Commands

```bash
# View logs
./scripts/docker-dev.sh logs
./scripts/docker-dev.sh logs foundation-app

# Execute commands in containers
./scripts/docker-dev.sh exec foundation-app go version
./scripts/docker-dev.sh exec foundation-db psql -U dev_user -d great_nigeria_foundation_dev

# Check status
./scripts/docker-dev.sh status

# Restart services
./scripts/docker-dev.sh restart

# Clean up everything
./scripts/docker-dev.sh clean
```

## Production Deployment

### 1. Prepare Environment

```bash
# Copy and configure production environment
cp .env.example .env.prod

# Edit production settings
nano .env.prod
```

### 2. Deploy with Production Configuration

```bash
# Deploy with production optimizations
docker-compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d

# Deploy with monitoring
docker-compose -f docker-compose.yml -f docker-compose.prod.yml --profile monitoring --env-file .env.prod up -d
```

### 3. Production Services

- **Application**: http://your-domain.com
- **Monitoring**: http://your-domain.com:3000 (Grafana)
- **Metrics**: http://your-domain.com:9090 (Prometheus)

## Configuration

### Environment Variables

Key environment variables for Docker deployment:

```bash
# Database
DB_PASSWORD=secure_password
DB_USERNAME=foundation_user
DB_DATABASE=great_nigeria_foundation

# Application
JWT_SECRET=your-32-character-secret
ENVIRONMENT=production
LOG_LEVEL=info

# External Ports
SERVER_EXTERNAL_PORT=8080
DB_EXTERNAL_PORT=5433
REDIS_EXTERNAL_PORT=6380
```

### Volume Management

Persistent data is stored in:

- **Database**: `foundation_db_data` volume
- **Redis**: `foundation_redis_data` volume
- **Uploads**: `./uploads` directory
- **Logs**: `./logs` directory

### Resource Limits

Default resource limits:

- **Database**: 512MB RAM, 0.5 CPU
- **Redis**: 256MB RAM, 0.25 CPU
- **Application**: 512MB RAM, 0.5 CPU
- **Nginx**: 128MB RAM, 0.25 CPU

## Monitoring and Logging

### Health Checks

All services include health checks:

- **Database**: `pg_isready` command
- **Redis**: `redis-cli ping` command
- **Application**: HTTP health endpoint
- **Nginx**: HTTP status check

### Logging

Logs are configured with:

- **Format**: JSON for structured logging
- **Rotation**: 10MB max size, 3 files retained
- **Location**: `/var/log/` in containers, `./logs/` on host

### Monitoring (Production)

Optional monitoring stack:

- **Prometheus**: Metrics collection
- **Grafana**: Visualization and dashboards
- **Health Checks**: Automated service monitoring

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check what's using the port
   lsof -i :8080
   
   # Change port in .env file
   SERVER_EXTERNAL_PORT=8081
   ```

2. **Permission Issues**
   ```bash
   # Fix upload directory permissions
   sudo chown -R 1001:1001 uploads/
   ```

3. **Database Connection Issues**
   ```bash
   # Check database logs
   docker-compose logs foundation-db
   
   # Test database connection
   docker-compose exec foundation-db psql -U foundation_user -d great_nigeria_foundation
   ```

4. **Memory Issues**
   ```bash
   # Check Docker memory usage
   docker stats
   
   # Increase Docker memory limit in Docker Desktop
   ```

### Debugging

```bash
# Enter application container
docker-compose exec foundation-app sh

# Check application logs
docker-compose logs -f foundation-app

# Monitor resource usage
docker stats

# Check health status
docker-compose ps
```

### Backup and Restore

```bash
# Backup database
docker-compose exec foundation-db pg_dump -U foundation_user great_nigeria_foundation > backup.sql

# Restore database
docker-compose exec -T foundation-db psql -U foundation_user great_nigeria_foundation < backup.sql

# Backup volumes
docker run --rm -v foundation_db_data:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz -C /data .
```

## Security Considerations

### Production Security

1. **Change Default Passwords**: Update all default passwords
2. **Use Strong Secrets**: Generate secure JWT secrets
3. **Network Security**: Use internal networks for service communication
4. **SSL/TLS**: Configure HTTPS with proper certificates
5. **Resource Limits**: Set appropriate CPU and memory limits
6. **Regular Updates**: Keep Docker images updated

### Development Security

1. **Isolated Environment**: Development uses separate databases
2. **Non-Root Users**: All containers run as non-root users
3. **Read-Only Mounts**: Configuration files mounted read-only
4. **Network Isolation**: Services communicate through Docker networks

## Performance Optimization

### Database Optimization

- Connection pooling configured
- Optimized PostgreSQL settings for workload
- Regular VACUUM and ANALYZE operations

### Application Optimization

- Multi-stage Docker builds for smaller images
- Proper resource limits and reservations
- Health checks for automatic recovery

### Caching Strategy

- Redis for session storage and caching
- Nginx for static file serving
- Proper cache headers for static assets

For more detailed information, see the [Configuration Guide](configuration.md) and [API Documentation](api-reference.md).
