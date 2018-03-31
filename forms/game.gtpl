{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        Game: {{ .Name }}<p>
        Number of Players: {{ .TotalPlayers }}<p>
        Round: {{ .IsNight }} {{ .Round }}<p>
        <b>Actions</b><br>
        {{ range $name, $num, $actions := .Players }}
            Player {{ $num }}: {{ $actions }}<br>
        {{ end }}    
        <p>
        <b>Players</b><br>
        {{ range $name, $num, $actions := .Players }}
            {{ $name }}<br>
        {{ end }}    
    </body>
</html>
{{end}}