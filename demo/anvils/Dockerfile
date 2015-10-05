# Since this is a Go project, we'll use the base golang image.
FROM golang

# Empire expects web processes to bind to the PORT environment variable, we'll set a sensible default in
ENV PORT=80

# Empire will extract the Procfile from the WORKDIR
WORKDIR /go/src/github.com/remind101/reinvent-2015/demo/anvils

# Add the source to the image and compile the app.
ADD . /go/src/github.com/remind101/reinvent-2015/demo/anvils
RUN go install github.com/remind101/reinvent-2015/demo/anvils

# We'll run the API command, and print the help by default.
CMD ["anvils", "-h"]
