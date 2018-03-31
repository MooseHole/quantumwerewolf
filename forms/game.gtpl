{{ define "selectContent" }}
    {{ range .PlayersByName }}
        <option value="{{ .Name }}">{{ .Name }}</option>
    {{ end }}   
{{ end }}

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
        <table>
        {{ range .PlayersByName }}        
            <tr>
            <form name={{ .Name }}>
            <td>{{ .Name }}</td>
            <td><select name="Attack">{{ template "selectContent" . }}</select></td>
            <td><select name="Peek">{{ template "selectContent" . }}</select></td>
            </tr>
            </form>
        {{ end }}   
     </body>
</html>
{{end}}