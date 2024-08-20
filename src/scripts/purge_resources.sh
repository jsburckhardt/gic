#!/bin/bash

# Check if user is logged in and wants to continue

read -p "You are about to purge resources. Did you login and select the correct subs, do you want to contibue? (y/n): " choice
if [[ $choice == "n" ]]; then
    exit 0
fi

# Echo notification working on subs
subscription=$(az account show --query "name" -o tsv)
echo "Working on subscription $subscription"




# Set the subscription

# Delete all Key Vaults
echo "the following Key Vaults will be deleted:"
az keyvault list --query "[].name" -o tsv
read -p "Do you want to continue? (y/n): " choice
if [[ $choice == "n" ]]; then
    exit 0
fi

for keyvault in $(az keyvault list-deleted --query "[].name" -o tsv)
do
    az keyvault purge --name $keyvault --no-wait
done



# Purge deleted cognitive services
echo "the following Cognitive Services will be deleted:"
az cognitiveservices account list-deleted --query "[].{name:name, location:location, id:id}" -o table
read -p "Do you want to continue? (y/n): " choice
if [[ $choice == "n" ]]; then
    exit 0
fi
items=$(az cognitiveservices account list-deleted --query "[].{name:name, location:location, id:id}" -o json)
count=$(echo $items | jq '. | length')
for i in $(seq 0 $(($count - 1)))
do
    # get the item
    objec=$(echo $items | jq ".[$i]")
    # echo $objec
    name=$(echo $objec | jq -r '.name')
    location=$(echo $objec | jq -r '.location')
    id=$(echo $objec | jq -r '.id')
    resource_group=$(echo $id | cut -d'/' -f9)
    echo deleting $name in $location in $resource_group
    az cognitiveservices account purge --name $name --location $location --resource-group $resource_group --verbose
done
