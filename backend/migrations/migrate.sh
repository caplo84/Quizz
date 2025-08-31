#!/bin/bash
# filepath: scripts/migrate.sh

# Database Migration Automation Script
# Usage: ./scripts/migrate.sh [command] [version]
# Commands: up, down, reset, status, force

set -e  # Exit on any error

# Configuration
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-quiz_db_dev}"
DB_USER="${DB_USER:-quiz_user}"
DB_PASSWORD="${DB_PASSWORD:-dev_password}"

# Auto-detect migrations path
if [ -z "$MIGRATIONS_PATH" ]; then
    if [ -f "001_create_topics_table.up.sql" ]; then
        # We're in the migrations directory
        MIGRATIONS_PATH="."
    elif [ -d "migrations" ] && [ -f "migrations/001_create_topics_table.up.sql" ]; then
        # We're in the parent directory
        MIGRATIONS_PATH="migrations"
    else
        # Default fallback
        MIGRATIONS_PATH="migrations"
    fi
fi

# Build database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if migrate tool is installed
check_migrate_tool() {
    if ! command -v migrate &> /dev/null; then
        log_error "migrate tool is not installed"
        log_info "Please install golang-migrate:"
        log_info "  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz"
        log_info "  sudo mv migrate /usr/local/bin/migrate"
        exit 1
    fi
}

# Test database connection
test_connection() {
    log_info "Testing database connection..."
    
    # First check if migration files exist
    if [ ! -d "$MIGRATIONS_PATH" ] || [ ! -f "$MIGRATIONS_PATH/001_create_topics_table.up.sql" ]; then
        log_error "Migration files not found in path: $MIGRATIONS_PATH"
        log_info "Please ensure you're running the script from the correct directory:"
        log_info "  - From project root: ./backend/migrations/migrate.sh"
        log_info "  - From migrations dir: ./migrate.sh"
        log_info "  - Or set MIGRATIONS_PATH environment variable"
        exit 1
    fi
    
    if migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" version &>/dev/null; then
        log_success "Database connection successful"
    else
        log_error "Cannot connect to database"
        log_info "Please check your database configuration:"
        log_info "  Host: $DB_HOST"
        log_info "  Port: $DB_PORT"
        log_info "  Database: $DB_NAME"
        log_info "  User: $DB_USER"
        log_info "  Migrations Path: $MIGRATIONS_PATH"
        exit 1
    fi
}

# Show help
show_help() {
    echo "Database Migration Tool"
    echo ""
    echo "Usage: $0 [command] [version]"
    echo ""
    echo "Commands:"
    echo "  up [N]          Apply all or N migrations"
    echo "  down [N]        Rollback all or N migrations"
    echo "  goto VERSION    Migrate to specific version"
    echo "  drop            Drop everything inside database"
    echo "  force VERSION   Set version but don't run migration"
    echo "  version         Print current migration version"
    echo "  status          Show migration status"
    echo "  reset           Drop and re-apply all migrations"
    echo "  help            Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST         Database host (default: postgres)"
    echo "  DB_PORT         Database port (default: 5432)"
    echo "  DB_NAME         Database name (default: quiz_db_dev)"
    echo "  DB_USER         Database user (default: quiz_user)"
    echo "  DB_PASSWORD     Database password (default: dev_password)"
    echo "  MIGRATIONS_PATH Migration files path (default: migrations)"
    echo ""
    echo "Examples:"
    echo "  $0 up           # Apply all pending migrations"
    echo "  $0 up 1         # Apply next 1 migration"
    echo "  $0 down 1       # Rollback 1 migration"
    echo "  $0 goto 3       # Migrate to version 3"
    echo "  $0 status       # Show current status"
    echo "  $0 reset        # Reset database"
}

# Apply migrations
migrate_up() {
    local steps=${1:-""}
    
    log_info "Applying migrations..."
    
    if [ -n "$steps" ]; then
        log_info "Applying $steps migration(s)"
        migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" up "$steps"
    else
        log_info "Applying all pending migrations"
        migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" up
    fi
    
    log_success "Migrations applied successfully"
}

# Rollback migrations
migrate_down() {
    local steps=${1:-""}
    
    log_warning "Rolling back migrations..."
    
    if [ -n "$steps" ]; then
        log_warning "Rolling back $steps migration(s)"
        migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" down "$steps"
    else
        log_warning "Rolling back all migrations"
        read -p "This will rollback ALL migrations. Are you sure? (y/N): " confirm
        if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
            migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" down -all
        else
            log_info "Rollback cancelled"
            exit 0
        fi
    fi
    
    log_success "Rollback completed successfully"
}

# Migrate to specific version
migrate_goto() {
    local version=$1
    
    if [ -z "$version" ]; then
        log_error "Version number required for goto command"
        exit 1
    fi
    
    log_info "Migrating to version $version..."
    migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" goto "$version"
    log_success "Migrated to version $version"
}

# Force version
force_version() {
    local version=$1
    
    if [ -z "$version" ]; then
        log_error "Version number required for force command"
        exit 1
    fi
    
    log_warning "Forcing version to $version..."
    migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" force "$version"
    log_success "Version forced to $version"
}

# Show current version
show_version() {
    log_info "Current migration version:"
    migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" version
}

# Show migration status
show_status() {
    log_info "Migration status:"
    echo "Database URL: $DB_URL"
    echo "Migrations Path: $MIGRATIONS_PATH"
    echo ""
    
    # Get current version
    current_version=$(migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" version 2>/dev/null || echo "0")
    echo "Current Version: $current_version"
    
    # List available migrations
    echo ""
    echo "Available migrations:"
    if [ -d "$MIGRATIONS_PATH" ]; then
        for file in "$MIGRATIONS_PATH"/*.up.sql; do
            if [ -f "$file" ]; then
                basename "$file" .up.sql
            fi
        done
    else
        log_warning "Migrations directory not found: $MIGRATIONS_PATH"
    fi
}

# Drop database
drop_database() {
    log_error "Dropping all database objects..."
    read -p "This will DROP ALL database objects. Are you sure? (y/N): " confirm
    if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
        migrate -database "$DB_URL" -path "$MIGRATIONS_PATH" drop
        log_success "Database dropped successfully"
    else
        log_info "Drop cancelled"
        exit 0
    fi
}

# Reset database (drop + migrate up)
reset_database() {
    log_warning "Resetting database (drop + migrate up)..."
    read -p "This will DROP ALL data and re-apply migrations. Are you sure? (y/N): " confirm
    if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
        drop_database
        migrate_up
        log_success "Database reset completed successfully"
    else
        log_info "Reset cancelled"
        exit 0
    fi
}

# Main script logic
main() {
    local command=${1:-"help"}
    local param=$2
    
    # Check prerequisites
    check_migrate_tool
    
    case $command in
        "up")
            test_connection
            migrate_up "$param"
            ;;
        "down")
            test_connection
            migrate_down "$param"
            ;;
        "goto")
            test_connection
            migrate_goto "$param"
            ;;
        "force")
            test_connection
            force_version "$param"
            ;;
        "version")
            test_connection
            show_version
            ;;
        "status")
            test_connection
            show_status
            ;;
        "drop")
            test_connection
            drop_database
            ;;
        "reset")
            test_connection
            reset_database
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Handle script interruption
trap 'log_error "Script interrupted"; exit 130' INT TERM

# Run main function with all arguments
main "$@"
