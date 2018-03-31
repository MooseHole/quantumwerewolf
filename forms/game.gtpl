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
{{template "PlayersAgain"}}
		{yield PLAYERSAGAINCONTENT}
{{end}}

        {{ end }}   

        {{ range .PlayersByName }}

		{{call PlayersAgain}}
			{{container PLAYERSAGAINCONTENT = .Name}}
        {{ end }}   
	
        {{ end }}   
     </body>
</html>
{{end}}