# go-discord-minesweeper

This is a simple Discord bot that generates Minesweeper boards.

# Usage

Download the binary or compile the source, then put your bot token in a
`config.json` file in the same directory as the bot:
```json
{
  "token": "<YOUR_BOT_TOKEN>"
}
```

Launch the bot, and type `!minesweeper` anywhere the bot has access to, and it
should spit out a perfectly workable 10x10 Minesweeper board with 10 mines.

# Advanced Usage

You can request a custom grid by typing `!minesweeper <width> <height>
<number-of-mines>`. It may refuse your request if the grid is too big, or if
more mines have been requested than cells (resulting in a literally unusable
grid).

You can set override the default grid settings in the `config.json` file:
```json
{
  "token": "<YOUR_BOT_TOKEN>",
  "defaultGridSize": [<WIDTH>, <HEIGHT>],
  "defaultMines": <NUMBEROFMINES>
}
```
