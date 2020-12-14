# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server

# create an image from apline
FROM alpine:latest

# copy over Go executable
COPY gateway gateway

# copy over handlers folder
COPY handlers handlers
COPY models models
COPY sessions sessions

# use the executable as the main entrypoint when the container starts
EXPOSE 443 80
ENTRYPOINT [ "/gateway" ]


