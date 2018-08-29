[![Go report](http://goreportcard.com/badge/kkragenbrink/slate)](http://goreportcard.com/report/kkragenbrink/slate) 
[![Build Status](https://travis-ci.com/kkragenbrink/slate.svg?branch=master)](https://travis-ci.com/kkragenbrink/slate)

# Slate
A [Discord Bot](http://www.discordapp.com/) written in Golang for managing Roleplaying Characters and Campaigns.

[Add Slate to your Server](https://discordapp.com/oauth2/authorize?client_id=484419646901059594&scope=bot)

## Features
- Create and list characters in discord
- View and edit character sheets from the website
- Roll dice from the character sheet to discord

## Deployment
Slate is deployed as a [heroku](http://www.heroku.com) application which hosts the SlateBot as well as the associated 
website, http://slate.sosly.org/.

## Data Storage and Security
All of Slate's data is stored in a heroku postgres cluster. Slate does not keep track of any information from Discord 
which is not documented, below.

### Discord Fields
Slate stores the following information from Discord in its database:

##### GuildID
Slate uses the GuildID to group together Channels which can be sent to. 

##### ChannelID
The ChannelID is tracked so that Slate knows which channels are configured to receive messages from Slate.

##### UserID
Your UserID is only used to determine which sheets you own.
