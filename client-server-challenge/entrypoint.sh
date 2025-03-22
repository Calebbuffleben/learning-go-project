#!/bin/bash
set -e

# Print system and environment information
echo "Container Startup Information:"
echo "-------------------------"
echo "Hostname: $(hostname)"
echo "Current User: $(whoami)"
echo "Working Directory: $(pwd)"
echo "Environment Variables:"
env | sort

# Execute the command passed to the container
if [ $# -gt 0 ]; then
    echo -e "\nExecuting command: $@"
    
    # Capture start time
    START_TIME=$(date +%s)
    
    # Run the command and capture output
    set +e  # Disable immediate exit on error
    "$@" 2>&1 | tee /app/command_output.log
    CMD_EXIT_CODE=${PIPESTATUS[0]}
    set -e  # Re-enable immediate exit on error
    
    # Calculate runtime
    END_TIME=$(date +%s)
    RUNTIME=$((END_TIME - START_TIME))
    
    echo -e "\nCommand Execution Details:"
    echo "Exit Code: $CMD_EXIT_CODE"
    echo "Runtime: $RUNTIME seconds"
    
    if [ $CMD_EXIT_CODE -ne 0 ]; then
        echo "ERROR: Command failed with exit code $CMD_EXIT_CODE"
        echo "Detailed output saved to /app/command_output.log"
    fi
fi

# Keep the container running for debugging
echo -e "\nKeeping container alive for debugging..."
tail -f /dev/null 