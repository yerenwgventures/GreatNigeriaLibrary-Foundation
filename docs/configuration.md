# Configuration Guide

This document explains how to configure the Great Nigeria Library Foundation application.

## Configuration Methods

The application supports multiple configuration methods with the following priority order:

1. **Environment Variables** (highest priority)
2. **YAML Configuration File** (medium priority)
3. **Default Values** (lowest priority)

## Backend Configuration

### YAML Configuration

Create a `config.yaml` file in the backend directory based on `config.example.yaml`:

```bash
cp backend/config.example.yaml backend/config.yaml
```

### Environment Variables

All YAML configuration values can be overridden using environment variables:

#### Database Configuration
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USERNAME` - Database username (default: postgres)
- `DB_PASSWORD` - Database password (required)
- `DB_DATABASE` - Database name (default: great_nigeria_foundation)
- `DB_SSL_MODE` - SSL mode (default: disable)

#### Authentication Configuration
- `JWT_SECRET` - JWT signing secret (required, minimum 32 characters)
- `ACCESS_TOKEN_EXPIRATION` - Access token expiration (default: 15m)
- `REFRESH_TOKEN_EXPIRATION` - Refresh token expiration (default: 168h)

#### Server Configuration
- `SERVER_HOST` - Server bind address (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/staging/production)

#### Redis Configuration
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (optional)
- `REDIS_DATABASE` - Redis database number (default: 0)

#### Email Configuration (Optional)
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `FROM_EMAIL` - From email address
- `FROM_NAME` - From name

#### OAuth Configuration (Optional)
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `GOOGLE_REDIRECT_URL` - Google OAuth redirect URL

#### Storage Configuration
- `STORAGE_TYPE` - Storage type (local/s3, default: local)
- `STORAGE_LOCAL_PATH` - Local storage path (default: ./uploads)
- `S3_REGION` - AWS S3 region
- `S3_BUCKET` - AWS S3 bucket name
- `S3_ACCESS_KEY` - AWS S3 access key
- `S3_SECRET_KEY` - AWS S3 secret key

#### Logging Configuration
- `LOG_LEVEL` - Log level (debug/info/warn/error, default: info)
- `LOG_FORMAT` - Log format (json/text, default: json)
- `LOG_OUTPUT` - Log output (stdout/file, default: stdout)
- `LOG_FILE` - Log file path (default: app.log)

## Frontend Configuration

### Environment Variables

Create a `.env.local` file in the frontend directory based on `.env.example`:

```bash
cp frontend/.env.example frontend/.env.local
```

#### API Configuration
- `REACT_APP_API_BASE_URL` - Backend API base URL (default: http://localhost:8080/api)
- `REACT_APP_API_TIMEOUT` - API request timeout in milliseconds (default: 30000)

#### Application Configuration
- `REACT_APP_APP_NAME` - Application name
- `REACT_APP_APP_VERSION` - Application version
- `REACT_APP_ENVIRONMENT` - Environment (development/staging/production)

#### Authentication Configuration
- `REACT_APP_JWT_STORAGE_KEY` - JWT token storage key
- `REACT_APP_REFRESH_TOKEN_KEY` - Refresh token storage key
- `REACT_APP_SESSION_TIMEOUT` - Session timeout in milliseconds

#### Feature Flags
- `REACT_APP_ENABLE_REGISTRATION` - Enable user registration
- `REACT_APP_ENABLE_OAUTH` - Enable OAuth authentication
- `REACT_APP_ENABLE_DISCUSSIONS` - Enable discussion features
- `REACT_APP_ENABLE_BOOKMARKS` - Enable bookmark features
- `REACT_APP_ENABLE_NOTES` - Enable note-taking features

## Docker Configuration

### Environment File

Create a `.env` file in the project root based on `.env.example`:

```bash
cp .env.example .env
```

This file contains all environment variables used by Docker Compose.

### Docker Compose Override

For local development customizations, create a `docker-compose.override.yml` file:

```yaml
version: '3.8'
services:
  foundation-app:
    volumes:
      - ./backend:/app/backend:ro
    environment:
      LOG_LEVEL: debug
```

## Security Considerations

### Required Secrets

The following configuration values must be changed from defaults in production:

1. **JWT_SECRET** - Use a cryptographically secure random string (minimum 32 characters)
2. **DB_PASSWORD** - Use a strong database password
3. **REDIS_PASSWORD** - Use a strong Redis password (if Redis auth is enabled)

### Example Secure Values

```bash
# Generate secure JWT secret
openssl rand -base64 32

# Generate secure passwords
openssl rand -base64 24
```

## Environment-Specific Configuration

### Development
- Use local database and Redis
- Enable debug logging
- Disable email verification
- Use mock OAuth credentials

### Staging
- Use managed database services
- Enable email notifications
- Use real OAuth credentials
- Enable rate limiting

### Production
- Use managed database services with SSL
- Enable all security features
- Use CDN for static assets
- Enable monitoring and alerting
- Use environment variables for all secrets

## Configuration Validation

The application validates configuration on startup and will fail to start if required values are missing or invalid.

Common validation errors:
- Missing JWT_SECRET
- Invalid database connection parameters
- Invalid email configuration (if email features are enabled)
- Invalid OAuth configuration (if OAuth is enabled)

## Troubleshooting

### Configuration Loading Issues

1. Check file permissions on config files
2. Verify YAML syntax is valid
3. Ensure environment variables are properly set
4. Check application logs for specific error messages

### Database Connection Issues

1. Verify database is running and accessible
2. Check database credentials
3. Verify network connectivity
4. Check SSL configuration

### Redis Connection Issues

1. Verify Redis is running and accessible
2. Check Redis password (if configured)
3. Verify network connectivity
4. Check Redis database number

For more help, see the troubleshooting section in the main README or contact support.
