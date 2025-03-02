## About
Telegram bot that allows you to track GitHub repositories and StackOverflow questions for the latest activity

## How to run

### Step 1: You need to create .env files in the root of the project

Create a `bot.env` file with the following content:

```
BOT_TOKEN=your_bot_token_here
SCRAPPER_URL=your_scrapper_url_here
SERVER_HOST=your_server_host_here
SERVER_PORT=your_server_port_here
```
Create a `scrapper.env` file with the following content:

```
SERVER_HOST=your_server_host_here
SERVER_PORT=your_server_port_here
BOT_URL=your_bot_url_here
```

### Step 2: Launch the Bot and Scraper

Use the Makefile to run the projects:

- To run the bot, use the command:
  ```
  make run_bot
  ```

- To run the scraper, use the command:
  ```
  make run_scrapper
  ```