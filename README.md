## SUI
*a startpage for your server and / or new tab page*

![screenshot](https://i.imgur.com/J4d7Q3D.png)

[More screenshots](https://imgur.com/a/FDVRIyw)

### Deploy with Docker compose

#### Prerequisites:
 - Docker: [Linux](https://docs.docker.com/install/linux/docker-ce/debian/), [Mac](https://hub.docker.com/editions/community/docker-ce-desktop-mac), [Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows)
 - [Docker-compose](https://docs.docker.com/compose/install/) 

### Shameless Plug

This was designed to work with my compose setup. This setup can be found [here](https://github.com/willfantom/composing). This is a good reference if you want too see how the configuration works using both docker labels and the traefik API.

### Configuration

This version of SUI is designed to pull the apps list from an external provider.

Currently supported providers are:
 - `Docker` via socket (`docker`)
 - `Træfik` via API (`traefik`)

### Core Configuration

(example [here](examples/config.json))

Copy the example in how bookmarks and search engines are added, it should be pretty simple 👍

For the providers, it is an array of `name`, `type` objects. The type must be one of the types above and the name must be another file in the same dir (excluding the `.json` bit).
```json
"appproviders": [
    { "name": "vps", "type": "traefik" },
    { "name": "main-traefik", "type": "traefik" }
  ]
```
If you're note using traefik, you can just use a docker provider by itself. 

#### Provider | Docker

This provider is the simplest to use. Add a file with the same name as you provided in the core config file, but add `.json` to it. Below is an example of a docker provider config that uses the local docker socket, provided you mount the socket to the container.

example:
```json
{
  "connection": "unix",
  "path": "/var/run/docker.sock"
}

```

You can then add flags to the containers inn the provided docker instance, such as:

example in compose:
```yaml
services:
   example:
      container_name: ex
      image: aservice:latest
      networks:
         - traefik-proxy
      labels:
         - [TRAEFIK LABELS]
         - sui.enabled=true
         - sui.icon=application
         - sui.name=Other Name
```

You must of course also mount the docker socket (as read-only).
Once this has been added, you can add labels to you containers (best via docker-compose files) that modify their sui behavior, for example:
 - `sui.enabled=true` will hide the service from the dashboard if false
 - `sui.name=Example` will set the application's name to `Example`
 - `sui.url=https://a.example.tld` will link the application to the given value
 - `sui.icon=application` will set the apps icon to `application` (see [here](https://materialdesignicons.com/))

#### Provider | Træfik

This provider uses the træfik API to determine what services are running. You can add a træfik provider to you SUI though the configuration JSON, for example:

```json
{
  "url": "https://traefik.mymainexample.tld",
  "user": "myusename",
  "pass": "supasecure1234",
  "ignored": ["PROMETHEUS@INTERNAL", "NOOP@INTERNAL", "API@INTERNAL"],
  "dockers": ["local-docker"]
}
```

- `url` specifies the URL of the target træfik instance
- `ignored` specifies service names to ignore from the apps list
- `dockers` list of docker type providers. These must be configured like the docker providers mentioned above
If basic auth is enabled on the Træfik endpoint:
- `user` and `pass` specify the credentials for basic auth

### Customization

#### Changing color themes
 - Click the options button on the left bottom

#### Color themes
These can be added or customized in the themer.js file. When changing the name of a theme or adding one, make sure to edit this section in index.html accordingly:

```
    <section  class="themes">
```

I might add a simpler way to edit themes at some point, but adding the current ones should be pretty straight forward.

### TODO

- [ ] Ignore traefik service with a regex
- [ ] Connect to remote docker instances via TCP
  - [ ] Ensure auth works with this too
- [ ] If service goes missing, perhaps gray out for a few refresh cycles
- [ ] Password protect `protected` services
- [ ] Add some other providers (including a simple File)
- [ ] Add more default icons