## About

This Telegram bot monitors the latest activity on GitHub repositories, such as issues and pull requests, as well as Stack Overflow questions, including answers and comments.

## Architecture

![Architecture Diagram](assets/architecture.png)

## How to Run

The bot can be launched using **Docker Compose**.

### Method 1: Using an `.env` File
1. Create a `.env` file in the root of the project with the following variables:
  ```
  BOT_TOKEN=<your_bot_token>
  POSTGRES_PASSWORD=<your_database_password>
  ```
2. Start the services using Docker Compose:
  ```
  docker-compose up -d
  ```

### Method 2: Passing Variables via Command Line
1. Pass the required variables directly in the command line:
  ```
  BOT_TOKEN=<your_bot_token> POSTGRES_PASSWORD=<your_database_password> docker-compose up -d
  ```
