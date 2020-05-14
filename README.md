## SUI
*a startpage for your server and / or new tab page*

![screenshot](https://i.imgur.com/J4d7Q3D.png)

[More screenshots](https://imgur.com/a/FDVRIyw)

### Deploy with Docker compose

#### Prerequisites:
 - Docker: [Linux](https://docs.docker.com/install/linux/docker-ce/debian/), [Mac](https://hub.docker.com/editions/community/docker-ce-desktop-mac), [Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows)
 - [Docker-compose](https://docs.docker.com/compose/install/) 

### Configuration

This version of SUI is designed to pull the apps list from an external provider.

Currently supported providers are:
 - Docker (via socket, tcp is a todo)
 - Træfik (via API) 

#### Provider | Docker

This provider is the simplest to use. Simply modify the config JSON adding the path and type. Currently only type `socket` is supported, but `TCP` will be added some day...

example:
```json
{
   ...
   "dockers": {
      "local": {
         "type": "socket",
         "path": "/var/run/docker.sock"
      }
   }
   ...
}
```

You must of course also mount the docker socket (as read-only).
Once this has been added, you can add labels to you containers (best via docker-compose files) that modify their sui behaviour, for example:
 - `sui.protected=true` will hide the service from the dashboard
 - `sui.name=Example` will set the application's name to `Example`
 - `sui.url=https://a.example.tld` will link the application to the given value
 - `sui.icon=application` will set the apps icon to `application` (see [here](https://materialdesignicons.com/))

#### Provider | Træfik

This provider uses the træfik API to determine what services are running. You can add a træfik provider to you SUI though the configuration JSON, for example:

```json
{
   ...
   "traefiks": {
      "example": {
         "url": "https://traefik.example.tld",
         "user": "myusername",
         "pass": "securatah",
         "ignored": "PROMETHEUS@INTERNAL NOOP@INTERNAL API@INTERNAL",
         "docker": "local"
      }
   }
   ...
}
```

- `url` specifies the URL of the target træfik instance
- `ignored` specifies service names to ignore from the apps list (separated by `' '`)
- `docker` can specify extra values from docker labels if have a docker provider also
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
