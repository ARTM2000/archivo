# Archivo
Archivo is created to be the final way of __archiving__ files, specially __configuration__ file or any type of documents on servers, by making file __backup management__ easy and useful in case of data loss or disaster occurrence.

- [Archivo](#archivo)
  - [How it works?](#how-it-works)
  - [How to run _archivo_ server?](#how-to-run-archivo-server)
  - [How to run _agent_?](#how-to-run-agent)


## How it works?
The design of archivo is on push mechanism which any server will push its desired files to one place to maintain snapshot of files by the time. There is two main components that act in this process, first `archivo server` (which we call `archivo` precisely) and second `archivo agent` (which we call `agent`).

`archivo` is responsible for file backups maintaining and file snapshot access management. Every servers that you want to send files from there to backup should be defined in `archivo` and they are called as `source server`. For each `source server` we have unique name and api key which will be use to authorized them on every file rotate.

`agent1`'s duty is to send file to `archivo`. It receives a configuration file that defines which file with what interval should be stored to `archivo` server and how many snapshot should it takes from that file.

> Note: In order to prevent any data loss, file rotate count decreasing is blocked. So you can only increase your files backup rotate count

## How to run _archivo_ server?
In order to run _archivo server_ you need first to download proper binary build of project that matches you target host (look at [release page](https://github.com/ARTM2000/archivo/releases) to download compressed binaries).

After you got the binary file, you need to have prepare your configuration file. By default, `archivo` tries to read configuration file from `${HOME}/.archivo.yaml`. You can find an example of `.archivo.yaml` [here](./example/server/.archivo.yaml).

`archivo` use postgres database to store its data, so you have to have it installed or run it with _docker_ or _docker-compose_.

After your configuration is ready, you should run `archivo` server by running:
```bash
# if you set your config file at ${HOME}/.archivo.yaml
./archivo

# if you set your config file elsewhere, pass the path of config file
./archivo -c /absolute/path/config/.archivo.yml
```

You can validate your configuration with following command:
```bash
# if you set your config file at ${HOME}/.archivo.yaml
./archivo validate

# if you set your config file elsewhere, pass the path of config file
./archivo validate -c /absolute/path/config/.archivo.yml
```

If everything was ok, your `archivo` server starts listening on `0.0.0.0:<PORT>` which PORT is your port number that you defined in config file. By default it starts listening on `8010`.

## How to run _agent_?
In order to run _agent server_ you need first to download proper binary build of project that matches you target host (look at [release page](https://github.com/ARTM2000/archivo/releases) to download compressed binaries).

After you got the binary file, you need to have prepare your configuration file. By default, `agent` tries to read configuration file from `${HOME}/.agent.yaml`. You can find an example of `.agent.yaml` [here](./example/agent/.agent.yaml).

After your configuration is ready, you should run `agent` by running:
```bash
# if you set your config file at ${HOME}/.agent1.yaml
./agent

# if you set your config file elsewhere, pass the path of config file
./agent -c /absolute/path/config/.agent.yml
```

You can validate your configuration with following command:
```bash
# if you set your config file at ${HOME}/.agent1.yaml
./agent validate

# if you set your config file elsewhere, pass the path of config file
./agent validate -c /absolute/path/config/.agent.yml
```

if everything goes successful, it will start to send files to archivo server
