#!/usr/bin/env bash

set -euo pipefail

#==============================================================================
# SQL Runner Script
# Provides unified interface to run SQL queries via Docker container or local psql
#==============================================================================

# Color codes for output
readonly RED='\033[0;31m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Global variables
ATTACH_MODE=true
VERBOSE=false
ENV_FILE=".env"
TEMP_SQL_FILE=""

#------------------------------------------------------------------------------
# Logging functions
#------------------------------------------------------------------------------
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*" >&2
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*" >&2
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_debug() {
    if [[ "$VERBOSE" == true ]]; then
        echo -e "${BLUE}[DEBUG]${NC} $*" >&2
    fi
}

#------------------------------------------------------------------------------
# Cleanup handler
#------------------------------------------------------------------------------
cleanup() {
    if [[ -n "$TEMP_SQL_FILE" && -f "$TEMP_SQL_FILE" ]]; then
        log_debug "Cleaning up temporary file: $TEMP_SQL_FILE"
        rm -f "$TEMP_SQL_FILE" || true
    fi
}

trap cleanup EXIT INT TERM

#------------------------------------------------------------------------------
# Help message
#------------------------------------------------------------------------------
show_help() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS] [SQL_QUERY_OR_FILE]

Run SQL queries via Docker container (default) or local psql.

OPTIONS:
    -a          Run inside Docker container (default)
    -c          Run using local psql
    -e FILE     Specify .env file path (default: .env in current directory)
    -v          Verbose/debug mode
    -h          Show this help message

ARGUMENTS:
    SQL_QUERY_OR_FILE
                SQL query string or path to .sql file
                If omitted, opens interactive psql session
                Can also read from stdin

EXAMPLES:
    $(basename "$0")                          # Interactive mode in container
    $(basename "$0") -c                       # Interactive mode locally
    $(basename "$0") "SELECT * FROM users;"   # Run query in container
    $(basename "$0") -c query.sql             # Run file locally
    cat file.sql | $(basename "$0")           # Read from stdin
    DRY_RUN=1 $(basename "$0") "SELECT 1;"    # Dry run mode

ENVIRONMENT:
    DRY_RUN=1   Print commands without executing

EOF
}

#------------------------------------------------------------------------------
# Load and validate .env file
#------------------------------------------------------------------------------
load_env() {
    if [[ ! -f "$ENV_FILE" ]]; then
        log_error ".env file not found: $ENV_FILE"
        log_error "Expected variables: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT, POSTGRES_CONTAINER_NAME"
        exit 1
    fi

    log_debug "Loading environment from: $ENV_FILE"

    # Source the .env file, handling comments and empty lines
    set -a
    while IFS= read -r line || [[ -n "$line" ]]; do
        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// }" ]] && continue
        
        # Remove 'export ' prefix if present and evaluate
        line="${line#export }"
        eval "$line" 2>/dev/null || true
    done < "$ENV_FILE"
    set +a

    # Validate required variables
    local required_vars=("POSTGRES_USER" "POSTGRES_PASSWORD" "POSTGRES_DB" "POSTGRES_PORT" "POSTGRES_CONTAINER_NAME")
    local missing_vars=()

    for var in "${required_vars[@]}"; do
        if [[ -z "${!var:-}" ]]; then
            missing_vars+=("$var")
        fi
    done

    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        log_error "Missing required environment variables: ${missing_vars[*]}"
        exit 1
    fi

    log_debug "Environment loaded successfully"
    if [[ "$VERBOSE" == true ]]; then
        log_debug "POSTGRES_USER=$POSTGRES_USER"
        log_debug "POSTGRES_DB=$POSTGRES_DB"
        log_debug "POSTGRES_PORT=$POSTGRES_PORT"
        log_debug "POSTGRES_CONTAINER_NAME=$POSTGRES_CONTAINER_NAME"
    fi
}

#------------------------------------------------------------------------------
# Test database connection
#------------------------------------------------------------------------------
test_connection() {
    local mode=$1
    log_debug "Testing database connection..."

    if [[ "$mode" == "attach" ]]; then
        if ! docker exec "$POSTGRES_CONTAINER_NAME" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1;" >/dev/null 2>&1; then
            log_error "Database connection test failed"
            return 1
        fi
    else
        if ! PGPASSWORD="$POSTGRES_PASSWORD" psql -h localhost -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1;" >/dev/null 2>&1; then
            log_error "Database connection test failed"
            return 1
        fi
    fi

    log_debug "Database connection successful"
    return 0
}

#------------------------------------------------------------------------------
# Attach mode: Run SQL inside Docker container
#------------------------------------------------------------------------------
run_attach_mode() {
    local sql_input="$1"

    log_info "Running in ATTACH mode (Docker container)"

    # Check if Docker is available
    if ! command -v docker >/dev/null 2>&1; then
        log_error "Docker command not found. Please install Docker."
        exit 1
    fi

    # Check if container exists
    if ! docker ps -a --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER_NAME}$"; then
        log_error "Container '$POSTGRES_CONTAINER_NAME' does not exist"
        exit 1
    fi

    # Check if container is running
    if ! docker ps --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER_NAME}$"; then
        log_error "Container '$POSTGRES_CONTAINER_NAME' is not running"
        exit 1
    fi

    # Check if psql is available in container
    if ! docker exec "$POSTGRES_CONTAINER_NAME" which psql >/dev/null 2>&1; then
        log_error "psql not found inside container '$POSTGRES_CONTAINER_NAME'"
        exit 1
    fi

    # Test connection
    test_connection "attach" || exit 1

    # Handle different input types
    if [[ -z "$sql_input" ]]; then
        # Interactive mode
        log_info "Starting interactive psql session..."
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB"
        else
            docker exec -it "$POSTGRES_CONTAINER_NAME" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
        fi
    elif [[ -f "$sql_input" ]]; then
        # File input
        log_info "Executing SQL file: $sql_input"
        
        local temp_container_path="/tmp/sql_script_$(date +%s)_$$.sql"
        
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: docker cp $sql_input $POSTGRES_CONTAINER_NAME:$temp_container_path"
            log_warn "DRY_RUN: docker exec $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB -f $temp_container_path"
            log_warn "DRY_RUN: docker exec $POSTGRES_CONTAINER_NAME rm -f $temp_container_path"
        else
            docker cp "$sql_input" "$POSTGRES_CONTAINER_NAME:$temp_container_path"
            
            # Ensure cleanup even on error
            trap "docker exec $POSTGRES_CONTAINER_NAME rm -f $temp_container_path 2>/dev/null || true; cleanup" EXIT INT TERM
            
            docker exec "$POSTGRES_CONTAINER_NAME" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$temp_container_path"
            docker exec "$POSTGRES_CONTAINER_NAME" rm -f "$temp_container_path"
        fi
    else
        # String query
        log_info "Executing SQL query"
        log_debug "Query: $sql_input"
        
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: docker exec $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB -c <SQL>"
        else
            docker exec "$POSTGRES_CONTAINER_NAME" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "$sql_input"
        fi
    fi
}

#------------------------------------------------------------------------------
# Connect mode: Run SQL using local psql
#------------------------------------------------------------------------------
run_connect_mode() {
    local sql_input="$1"

    log_info "Running in CONNECT mode (local psql)"

    # Check if local psql is available
    if ! command -v psql >/dev/null 2>&1; then
        log_error "psql command not found. Please install PostgreSQL client."
        exit 1
    fi

    # Test connection
    test_connection "connect" || exit 1

    # Prepare psql connection parameters
    export PGPASSWORD="$POSTGRES_PASSWORD"
    local psql_args=(-h localhost -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB")

    # Handle different input types
    if [[ -z "$sql_input" ]]; then
        # Interactive mode
        log_info "Starting interactive psql session..."
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB"
        else
            psql "${psql_args[@]}"
        fi
    elif [[ -f "$sql_input" ]]; then
        # File input
        log_info "Executing SQL file: $sql_input"
        
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f $sql_input"
        else
            psql "${psql_args[@]}" -f "$sql_input"
        fi
    else
        # String query
        log_info "Executing SQL query"
        log_debug "Query: $sql_input"
        
        if [[ "${DRY_RUN:-0}" == "1" ]]; then
            log_warn "DRY_RUN: psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c <SQL>"
        else
            psql "${psql_args[@]}" -c "$sql_input"
        fi
    fi
}

#------------------------------------------------------------------------------
# Main function
#------------------------------------------------------------------------------
main() {
    local sql_input=""

    # Parse options
    while getopts "ace:vh" opt; do
        case "$opt" in
            a)
                ATTACH_MODE=true
                ;;
            c)
                ATTACH_MODE=false
                ;;
            e)
                ENV_FILE="$OPTARG"
                ;;
            v)
                VERBOSE=true
                ;;
            h)
                show_help
                exit 0
                ;;
            *)
                log_error "Invalid option: -$OPTARG"
                show_help
                exit 1
                ;;
        esac
    done

    shift $((OPTIND - 1))

    # Load environment
    load_env

    # Determine SQL input source
    if [[ $# -gt 0 ]]; then
        # Command line argument
        sql_input="$*"
    elif [[ ! -t 0 ]]; then
        # Stdin input
        log_debug "Reading SQL from stdin"
        TEMP_SQL_FILE=$(mktemp)
        cat > "$TEMP_SQL_FILE"
        sql_input="$TEMP_SQL_FILE"
    fi

    # Run appropriate mode
    if [[ "$ATTACH_MODE" == true ]]; then
        run_attach_mode "$sql_input"
    else
        run_connect_mode "$sql_input"
    fi
}

# Entry point
main "$@"
