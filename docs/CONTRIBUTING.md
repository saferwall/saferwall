# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change. 

Please note we have a code of conduct, please follow it in all your interactions with the project.

### Repostiory Layout
* __pkg__ : library code to use by external applications.
* __website__ : saferwall website and documentation.
* __ui__ : (frontend) vue.js dashboard.
* __web__ : (backend) go web application.
* __build__ : docker files, makefiles, kubernetes deployements.
* __api__ : proto buffer specs.

### Requirements

```
- sudo apt-get install make curl
- touch .env
```

Once you do that, you can do nearly everything using make, you can look at the supported options by executing a `make`.
