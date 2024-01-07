#!/bin/bash

# Define environments
environments=("sandbox" "development" "staging" "production")

# Check for arguments
if [ "$#" -eq 1 ]; then
    operation=$1
    selected_environments=("${environments[@]}")
elif [ "$#" -eq 2 ]; then
    operation=$1
    selected_environments=("$2")
else
    echo "Usage: $0 [-e|-d] [optional: sandbox|development|staging|production]"
    exit 1
fi
# Function to perform sops operation
perform_sops_operation() {
    local environment=$1
    decrypted_file="fixtures.${environment}.tfvars.json"
    encrypted_file="fixtures.${environment}.tfvars.enc.json"

# Perform sops operation and check if successful
if [ "$operation" == "-e" ]; then
    if [ -f "$encrypted_file" ]; then
        echo "Encrypted file already exists. Skipping encryption."
        exit 0
    fi
    if sops -e "$decrypted_file"  > "$encrypted_file"; then
        echo "Encryption successful."
        rm -f "$decrypted_file"
        git add "$encrypted_file"
        echo "Original file deleted."
    else
        echo "Encryption failed."
        rm -f "$encrypted_file"
        exit 1
    fi
elif [ "$operation" == "-d" ]; then
    if [ -f "$decrypted_file" ]; then
        echo "Decrypted file already exists. Skipping decryption."
        exit 0
    fi
    if sops --ignore-mac -d "$encrypted_file" > "$decrypted_file"; then
        echo "Decryption successful."
        rm -f "$encrypted_file"
        echo "Encrypted file deleted."
    else
        echo "Decryption failed."
        rm -f "$decrypted_file"
        exit 1
    fi
else
    echo "Invalid operation: $operation"
    exit 1
fi
}

# Loop over selected environments
for env in "${selected_environments[@]}"; do
    perform_sops_operation "$env"
done
