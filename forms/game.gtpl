{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        Game: {{ .Name }}<br>
        Number of Players: {{ .TotalPlayers }}<br>
        Round: {{ .Round }}<p>
        <b>Actions</b><br>
        {{ range .PlayersByNum }}
            Player {{ .Num }} ({{ .Name }}): {{ .Actions }}<br>
        {{ end }}    
        <p>
        <b>Action Messages</b><br>
        {{ .ActionMessages }}
        <p>
        <b>Actions for {{ .Round }}</b><br>
        <form name="Actions" id="Actions" action="/processActions" method="post">
        <input type="hidden" name="gameId" value={{ .GameID }}>
        <input type="submit">
        <input type="checkbox" name="advance" value="true">Advance to next round<br>
        <table name="ActionsTable" id="ActionsTable">
        <tr><th>Player</th>{{ if .IsNight }}<th>Attack</th><th>Peek</th>{{ else }}<th>Lynch</th>{{ end }}</tr>
        {{ range $name, $selections := .ActionSubjects }}
        <tr>
            <td>{{ $name }}</td>
            {{ if $.IsNight }}
                <td>
                <select name="{{ $name }}Attack">
                {{ range $name, $value :=  $selections.Attack }}
                <option value="{{ $value }}">{{ $name }}</option>
                {{ end }}
                </select>
                </td><td>
                <select name="{{ $name }}Peek">
                {{ range $name, $value :=  $selections.Peek }}
                <option value="{{ $value }}">{{ $name }}</option>
                {{ end }}
                </select>
                </td>
            {{ else }}
                <td>
                <select name="{{ $name }}Lynch">
                {{ range $name, $value :=  $selections.Lynch }}
                <option value="{{ $value }}">{{ $name }}</option>
                {{ end }}
                </select>
                </td>
            {{ end }}
        </tr>
        {{ end }}
        </table>
        </form>
     </body>
</html>
{{end}}