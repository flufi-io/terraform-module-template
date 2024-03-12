#!/bin/bash
# This script is used to encrypt and decrypt secrets for different environments and regions.
# It uses sops for encryption and decryption.

# Define environments and regions
# An array of environments and regions for which the secrets need to be managed.
environments=("sandbox" "development" "staging" "production")
regions=("us-east-1" "us-west-2")

# Check for arguments
# The script expects one to three arguments.
# The first argument is the operation (-e for encryption, -d for decryption).
# The second argument is optional and specifies the environment. If not provided, the operation is performed for all environments.
# The third argument is optional and specifies the region. If not provided, the default region is 'us-east-1'.
if [ "$#" -eq 1 ]; then
    operation=$1
    selected_environments=("${environments[@]}")
    region="us-east-1"
elif [ "$#" -eq 2 ]; then
    operation=$1
    selected_environments=("$2")
    region="us-east-1"
elif [ "$#" -eq 3 ]; then
    operation=$1
    selected_environments=("$2")
    region=$3
else
    echo "Usage: $0 [-e|-d] [optional: ${regions[*]}] [optional: ${regions[*]}]"
    exit 1
fi

# Function to perform sops operation
# This function performs the sops operation (encryption or decryption) for a given environment and region.
perform_sops_operation() {
    local environment=$1
    local region=$2
    decrypted_file="secrets.${environment}.${region}.auto.tfvars.json"
    encrypted_file="secrets.${environment}.${region}.auto.tfvars.enc.json"

    # Perform sops operation and check if successful
    # If the operation is encryption, the script checks if the encrypted file already exists.
    # If it does, the script exits. If it doesn't, the script attempts to encrypt the decrypted file.
    # If the operation is decryption, the script checks if the decrypted file already exists.
    # If it does, the script exits. If it doesn't, the script attempts to decrypt the encrypted file.
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
# The script performs the sops operation for each selected environment.
for env in "${selected_environments[@]}"; do
    perform_sops_operation "$env" "$region"
done
