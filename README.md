# ledger
A full-stack web app that simplifies your budgeting routine.

## Setup
I recommend running ledger as a containerized application.  Please note that these instructions were tested on a Windows 11 machine running Ubuntu 24.04.2 on WSL.

1. Install Docker.  I recommend downloading [Docker Desktop](https://www.docker.com/products/docker-desktop/) for your machine's platform.
2. Open a terminal and run the following command from the root of the project (it may take a few minutes to complete):

   ```Bash
   docker compose -f 'compose.yml' up -d --build
   ```
   Once the process is complete, there will be 2 running containers (one for the Postgres service, and one for ledger).
3. Open a web browswer and navigate to ```localhost::8080```.  You should now see the ledger homepage!
   
## Usage
TBD