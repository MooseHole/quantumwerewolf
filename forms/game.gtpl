{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        Game: {{ .Name }}<P>
        Number of Players: {{ .TotalPlayers }}<P>
        Round: {{ .IsNight }} {{ .Round }}
    </body>
</html>
{{end}}