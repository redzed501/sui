# Configuring SUI

If all you need is a clean looking dashboard in which you configure all app information into a single simple JSON, look no further than SUI's upstream [here]()

## Basic Configuration

You should create a new file called `config.json`, or simply copy the one in [here](./config.json). 

Fields:
 - `title` - What is the title of you dash (in tab bar)
 - `debug` - Enable debug level logging
 - `app_refresh` - How many seconds between refreshing each app provider
 - `bookmakrs` - an object of named bookmark categories. Each category should have a list of bookmark objects. Each bookmark object should have its `name` and its `url`.
 - `engines` - an object of named search engine objects. Each search engine object should have its `url` and a `prefix` (used for quick access in the search bar)
 - `appproviders` - see below:

## App Providers

The name and type of each provider is configured in the `config.json` created above. such as:

```json
"appproviders": [
    { "name": "local-docker", "type": "docker" },
    { "name": "remote-traefik", "type": "traefik" },
    { "name": "random-traefik", "type": "traefik" }
  ],
```
each name should have a corresponding json file in the same directory as `config.json`, such as `random-traefik.json`.

Type can be any of the following:
 - `traefik`
 - `docker`

### Docker Provider

#### Unix Connection (docker.sock)

To use a local docker socket, ensure it is mounted to the container!

The config then looks like so:
```json
{
  "connection": "unix",
  "path": "/var/run/docker.sock"
}
```

#### TCP Connection

Example config:
```json
{
  "connection": "tcp",
  "url": "10.30.65.123:2375"
}
```

### Traefik Provider

Example config:
```json
{
  "url": "https://traefik.myrandomexample.tld",
  "user": "myusename",
  "pass": "supasecure1234",
  "ignored": ["PROMETHEUS@INTERNAL", "NOOP@INTERNAL", "API@INTERNAL"],
  "dockers": ["remote-docker"]
}
```

Any `dockers` in a traefik config do not have to have been in the `config.json`, just so long as they have a `<name>.json` file, it should work.

For basic auth, `user` and `pass` can be added.
