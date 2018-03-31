{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        Game: {{ .Name }}<br>
        Number of Players: {{ .TotalPlayers }}<br>
        Round: {{ .IsNight }} {{ .Round }}<p>
        <b>Actions</b><br>
        {{ range .PlayersByNum }}
            Player {{ .Num }}: {{ .Actions }}<br>
        {{ end }}    
        <p>
        <b>Players</b><br>
        {{ range .PlayersByName }}
            {{ .Name }}<br>
        {{ end }}    
    </body>
</html>
{{end}}