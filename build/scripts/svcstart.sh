#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
/svc/main &

# Start the helper process
ollama run llama3.2:1b

# now we bring the primary process back into the foreground
# and leave it there
fg %1
