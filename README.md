# ledger
A full-stack web app that simplifies your budgeting routine.

## Setup
I recommend running ledger as a containerized application, as it is not yet fully deployed.  Please note that these instructions were tested on a Windows 11 machine running Ubuntu 24.04.2 on WSL.

1. Clone this repo onto your local machine.
2. Install Docker.  I recommend downloading [Docker Desktop](https://www.docker.com/products/docker-desktop/) for your machine's platform.
3. Open a terminal and run the following command from the root of the project (it may take a few minutes to complete):

   ```Bash
   docker compose up
   ```
   Once the process is complete, 2 containers will be running on your machine (one for the Postgres service, and one for the web server).
4. Open a web browser and navigate to ```localhost::8080```.  You should now see the ledger homepage!
   
## Usage
TBD

## Troubleshooting
- Sometimes you may encounter an issue that has been fixed in a recent image.  To ensure you are using the most recent image when running via Docker compose, delete all existing ledger-related containers and images on your local machine first.