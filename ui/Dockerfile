################################
# STEP 1 build executable binary
################################

FROM node:10-alpine as build-stage

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# make the 'app' folder the current working directory
WORKDIR /app

# copy both 'package.json' and 'package-lock.json' (if available)
COPY package.json ./
COPY yarn.lock ./

# install project dependencies
RUN yarn install

# copy project files and folders to the current working directory (i.e. 'app' folder)
COPY . .
RUN yarn build


############################
# STEP 2 build a small image
############################

FROM nginx:stable-alpine as production-stage
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Saferwall UI Dashboard"

COPY --from=build-stage /app/dist /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]