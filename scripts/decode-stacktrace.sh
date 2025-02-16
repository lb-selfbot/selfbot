#!/bin/bash

# Prompt for stacktrace input
echo -n "Input stacktrace (CTRL+D to end): "

# Read stacktrace lines until EOF (Ctrl+D)
stacktrace=""
while IFS= read -r line; do
  stacktrace+="$line"$'\n'
done

# Write stacktrace to file
echo "$stacktrace" > stack.txt

# Run garble command
garble -seed=U4Nxq2YSVm6ivvOejQbElA== reverse -tags desktop github.com/LightningDev1/LB-Selfbot stack.txt

# Cleanup stack.txt
rm stack.txt