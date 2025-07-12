#!/bin/bash
# Great Nigeria Library Foundation - Development Docker Management Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating from .env.example..."
        if [ -f .env.example ]; then
            cp .env.example .env
            print_success "Created .env file from .env.example"
            print_warning "Please review and update the .env file with your settings"
        else
            print_error ".env.example file not found. Please create .env file manually."
            exit 1
        fi
    fi
}

# Function to start development environment
start_dev() {
    print_status "Starting Great Nigeria Library Foundation development environment..."
    
    check_docker
    check_env_file
    
    # Create necessary directories
    mkdir -p dev-uploads dev-logs data/postgres data/redis
    
    # Start services
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
    
    print_success "Development environment started!"
    print_status "Services available at:"
    echo "  - Application: http://localhost:8081"
    echo "  - Database Admin: http://localhost:8082 (if dev-tools profile enabled)"
    echo "  - Redis Admin: http://localhost:8083 (if dev-tools profile enabled)"
    echo "  - Mail Catcher: http://localhost:8025 (if dev-tools profile enabled)"
}

# Function to start with dev tools
start_dev_tools() {
    print_status "Starting development environment with dev tools..."
    
    check_docker
    check_env_file
    
    mkdir -p dev-uploads dev-logs data/postgres data/redis
    
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml --profile dev-tools up -d
    
    print_success "Development environment with dev tools started!"
    print_status "Services available at:"
    echo "  - Application: http://localhost:8081"
    echo "  - Database Admin: http://localhost:8082"
    echo "  - Redis Admin: http://localhost:8083"
    echo "  - Mail Catcher: http://localhost:8025"
}

# Function to stop development environment
stop_dev() {
    print_status "Stopping development environment..."
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
    print_success "Development environment stopped!"
}

# Function to restart development environment
restart_dev() {
    print_status "Restarting development environment..."
    stop_dev
    start_dev
}

# Function to view logs
logs_dev() {
    if [ -n "$1" ]; then
        docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f "$1"
    else
        docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f
    fi
}

# Function to execute commands in containers
exec_dev() {
    if [ -z "$1" ]; then
        print_error "Please specify a service name"
        exit 1
    fi
    
    service="$1"
    shift
    
    if [ $# -eq 0 ]; then
        docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec "$service" sh
    else
        docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec "$service" "$@"
    fi
}

# Function to show status
status_dev() {
    print_status "Development environment status:"
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps
}

# Function to clean up
clean_dev() {
    print_warning "This will remove all containers, networks, and volumes for the development environment."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Cleaning up development environment..."
        docker-compose -f docker-compose.yml -f docker-compose.dev.yml down -v --remove-orphans
        docker system prune -f
        print_success "Development environment cleaned up!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Function to show help
show_help() {
    echo "Great Nigeria Library Foundation - Development Docker Management"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start       Start development environment"
    echo "  start-tools Start development environment with dev tools"
    echo "  stop        Stop development environment"
    echo "  restart     Restart development environment"
    echo "  logs [svc]  Show logs (optionally for specific service)"
    echo "  exec <svc>  Execute command in service container"
    echo "  status      Show status of services"
    echo "  clean       Clean up all containers and volumes"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 logs foundation-app"
    echo "  $0 exec foundation-app go version"
}

# Main script logic
case "${1:-help}" in
    start)
        start_dev
        ;;
    start-tools)
        start_dev_tools
        ;;
    stop)
        stop_dev
        ;;
    restart)
        restart_dev
        ;;
    logs)
        logs_dev "$2"
        ;;
    exec)
        shift
        exec_dev "$@"
        ;;
    status)
        status_dev
        ;;
    clean)
        clean_dev
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
