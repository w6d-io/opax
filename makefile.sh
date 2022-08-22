#!/bin/bash

############################################################
# Help                                                     #
############################################################
Help()
{
   # Display Help
   echo "Description of the script functions here"
   echo
   echo "This script create configuration file for OPA and run a opa server "
   echo
   echo "options:"
   echo " - first   parameter \${1}  is the target action."
   echo "                - The basic action is <<config>> : use to create all the config files necessary for run a OPA"
   echo "                - The advanced action run a server OPA."
   echo "                      - <<run>>"
   echo
}

############################################################
############################################################
# Main program                                             #
############################################################
############################################################

# Display the description of the script functions here.
if [ $# -eq 0 ]
then
  Help
fi

# Create OPA config files if $1 is passed to this script with value "config"
if [[ "${1}" = "config" ]]; then
cat << EOF1 >  bin/main.rego
package system
main = true
str = "foo"
loopback = input
EOF1

cat << EOF2 >  bin/input.json
{"input":{"foo": "bar"}}
EOF2
fi

# Run opa  $1 is passed to this script with value "run"
if [[ "${1}" = "run" ]]; then

./bin/opa run --server ./bin/main.rego

fi
