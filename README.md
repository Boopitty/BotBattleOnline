# Bot Battle Online

This is my capstone project for boot.dev! 
Bot Battle is a turn based browser game, where both players choose a team and fight to last bot.

## Motivation

I wanted to make a project that was easily accessible to my family and potential employers.
A browser game made the most sense, because I already made one as an earlier project.
I'll have to re-write the code for the game in `go` so it will work more reliably with the server code, but the principle is mostly the same.

## The Game

### Team Building
Choose from a variety of bots with unique abilities to create your team.
Each bot has a point value, and you team can only be worth a certain amount.
Choose wisely, because a single bot could take up your entire team's value!

### Battle
Two bots can be active at a time, but some bots can be HUGE and take up both active slots.
Each bot makes one action per turn, even if they're HUGE, and in order of fastest to slowest.
Certain moves may have priority, which forces them to move faster than normal.

## Quick Start
### Prerequisites

- [Go](https://golang.org/doc/install) 1.26.1 or later
- [PostgreSQL](https://www.postgresql.org/download/)
- [goose](https://goose-docs.ai) used for database migrations
- [Node.js] 20+ and [npm] used for development

### Install

Clone the project and change directory.
```bash
git clone github.com/Boopitty/BotBattleOnline.git
cd BotBattleOnline
```

Download dependancies.
```bash
# Go dependancies
go mod download

# Change to bot-battle directory
cd bot-battle

# download JavaScript dependencies
npm install

# Return to root
cd ..
```

## Configuration
### Database Setup

Create new user and a database. 
Replace `username` with a new username, `password` with a new secure password, and `db_name` with whatever you want to name your database.

```bash
# Connect to postgres:
psql postgres 

# create your role:
CREATE USER username WITH PASSWORD 'password'; 

# Create the database:
CREATE DATABASE db_name OWNER username; 

# quit:
\q 

# Your connection string will be:
"postgres://username:password@localhost:5432/db_name"
```

Make an up migration to the database.
```bash
goose -dir sql/schema postgres "connection_string" up
```

### .env File

A .env.example file has been provided. Create a new `.env` file using it.
```bash
cp .env.example .env
```

Generate a new `SECRET` variable by running this command in your terminal. 

Copy it into the `.env` file once made.
```bash
openssl rand -base64 32
```

Save your database's connection string into your `.env` file as the `DB_URL` variable:
```bash
# Your connection string will be:
"postgres://username:password@localhost:5432/db_name"
```

## Startup

The frontend and backend are started separately in different terminals. 
There is a shell command file for each. Use `ctrl + c` to cancel them.

### Server Startup

Start the server first to avoid websocket issues.
```bash
# Builds and runs backend server
bash server.sh
```

### Frontend Startup

Start the frontend in a new terminal. Visit the given link.
```bash
# Build and run the frontend
bash bot-battle.sh
```

## Contributing

Thanks for having interest in my project! It's still a work in progress, but feel free to report any bugs and errors you find!