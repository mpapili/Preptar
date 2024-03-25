#! /bin/bash

# little helper tool for developing

docker run -it --rm -v "$(pwd):/app:Z" \
	--entrypoint /bin/bash \
	golang

