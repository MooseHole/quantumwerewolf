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
            {{ range _ctx.PlayersByNum }}
                Player {{ .Num }}: {{ .Actions }}<br>
            {{ end }}    
        {{ end }}    
     </body>
</html>
{{end}}