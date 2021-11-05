# Go Scripts

## Table of contents
* [Setup](#setup)

## Setup

### Build Executable

1. Use go package manager to get go-scripts: 
```
git clone https://github.com/joshwi/go-scripts.git
```

2. Change directory into repo
```
cd go-scripts
```
3. Create .env
```
nano .env
```
4. Paste environment variables in file
```
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=password123
NEO4J_SERVICE_HOST=localhost
NEO4J_SERVICE_PORT=7687
```
5. Build go code into executable
```
go build -o collector
```
6. Run the collector

Example: Get all games from 2020 NFL Season
```
./collector -c='pfr_map_team'
./collector -c='pfr_map_season'
./collector -c='pfr_team_season' -q='MATCH (n:pfr_map_team_teams),(m:pfr_map_season_years) WHERE m.year="2020" RETURN DISTINCT n.tag as tag, m.year as year'
```

### Build Docker Image

1. Use go package manager to get go-scripts: 
```
git clone https://github.com/joshwi/go-scripts.git
```

2. Change directory into repo
```
cd go-scripts
```
3. Create .env
```
nano .env
```
4. Paste environment variables in file
```
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=password123
NEO4J_SERVICE_HOST=localhost
NEO4J_SERVICE_PORT=7687
```
5. Build app in docker container image: 
```
sudo docker build -t collector .
```
6. Run docker container:
```
sudo docker run -d --net <docker_network> --env-file <path_to_env>  --name collector collector
```

### Publish Docker Image

1. Build app in docker container image: 
```
sudo docker build -t <container_repo>/collector .
```

2. Push docker image to container repo: 
```
sudo docker push <container_repo>/collector:latest
```