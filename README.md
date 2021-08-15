# dots-server
Server side application for [dots-cli](https://github.com/alvanrahimli/dots-cli)  
If you don't know what `dots-cli` is, you can follow [this link](https://github.com/alvanrahimli/dots-cli)
___

## Installation
dots-server is available as binary and docker container.
It is recommended to use provided `docker-compose.yml` file, as it is more up to date and easy to set up.  

- Install docker by referring [this url](https://docs.docker.com/engine/install/)
- Do post-install set-up from [this url](https://docs.docker.com/engine/install/linux-postinstall/)
- Install docker-compose from [this url](https://docs.docker.com/compose/install/)
- Copy `docker-compose.yml` from [here](#docker-compose-file) to the place you want to install `dots-server`
- Run `docker-compose up -d`. It will run server and expose port `9090`.
- It is recommended to use some kind of reverse proxy (e.g Nginx).
  - Configure nginx to forward requests to port: `9090`. 
  - Set `proxy_pass_request_headers	on;` to ensure that nginx will pass request headers to app.

```
docker-compose up -d      # Runs application
docker-compose logs -f    # Shows application logs and follows it.
docker-compose down       # Stops application
docker-compose kill       # Kills container
```

## Requirements
`dots-server` requires following configurations:
- Environemnt variables:
  - PORT: which port does application listen
  - DB_PATH: Sqlite database file path relative to dots-server executable
  - REGISTRY_DOMAIN: Domain that app is hosted
- Files & Directories
  - `archives` directory, same place with dots-server executable
  - Sqlite database file. Database creation SQL script can be found [here](db-template.sql) Path can be set via `DB_PATH` env. variable

## Docker Compose file
```
version: "3.9"
services:
  dots-server:
    image: dotsorg/dots-server
    ports:
      - "9090:80"
    environment:
      - PORT=80
      - REGISTRY_DOMAIN=your-awesome-domain.com
      - DB_PATH=dots.db
    volumes:
      - '/var/www/dots-server/archives:/app/archives/'
      - '/var/www/dots-server/dots.db:/app/dots.db'
```

## TODO (API):
- [x] Login / Register API endpoints
- [x] Push package API endpoint
- [x] Get Package Archive endpoint
- [ ] Update & Delete packages endpoint
- [ ] Enhanced models (info, settings etc.)
- [ ] Enhanced endpoints for webapp
- [ ] Security considerations
