# Telegram Bot for reminders

## What can this bot do?
> -***Hey bot, what can you do? What were you made for?***<br>
-*Hey, user! Thank you for your question. Let me explain you. I'm designed to help you not to forget about the things you should do. For example, if your friend asked you to pick him up at 7 p.m., you can just tell me about it  - and I will let you know in due time. I have a list of commands, such as `/add_reminder`, `/add_note`, etc.*<br>
-***Well done! Do you have any other functions?***<br>
-*Of course! You can also use me like a notebook: I will keep saved your notes as long as you need.*<br>

This bot can send you reminders about the things you should do. All you need to do is just to type the command `/add_reminder` (or choose it in the list of commands), then type the title, choose the date and time. That's all! Bot will notify you when  the time comes.
But this is not all! You can also use `/add_note` command to save note if you need. In the future bot will be able to save photos!

## How to install this bot?
> -***Bot, could you tell me, how can I install you on my local PC?***<br>
-*Sure. All you need is just to run command: `git clone github.com/Pizhlo/tgbot-reminder-go.git`. Next you should set the configuration: create file `app.env` in the root folder and define the following variables: your telegram token, which you got from Bot Father; and db source like `postgresql://user:pass@localhost:port/db_name?sslmode=disable`. That's all!*

To use this bot, first you need to get token from Bot Father.<br>
Then clone this repo:<br>
`git clone [github.com/Pizhlo/tgbot-reminder-go.git](https://github.com/Pizhlo/bot-reminder-go-telegram)`<br>
create `app.env` file (see `app.env.example`) and paste there your token and db source.
