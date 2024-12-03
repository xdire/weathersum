FROM golang:1.22.9 AS BUILDER

WORKDIR /svc

COPY . .
RUN go build -o main .

# COPY my_main_process my_main_process
# COPY my_helper_process my_helper_process
# COPY my_wrapper_script.sh my_wrapper_script.sh

FROM ollama/ollama AS APPLICATION

WORKDIR /svc
# Copy main program for the weather API pull
COPY --from=BUILDER /svc/main ./
COPY --from=BUILDER /svc/build/scripts ./

# Embed ollama model into the image
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
    curl
RUN ollama serve & \
    # sleep 5 \
    curl --retry 10 --retry-all-errors --connect-timeout 20 --retry-connrefused --retry-delay 5 http://localhost:11434/ && \
    curl -X POST -d '{"name": "llama3.2:1b"}' http://localhost:11434/api/pull

ENTRYPOINT ["/svc/svcstart.sh"]