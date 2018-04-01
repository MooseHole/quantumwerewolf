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
        <tr><th>Player</th><th>Attack</th><th>Peek</th></tr>
        {{ range .PlayersByName }}        
        <tr>
        <form name={{ .Name }} id={{ .Name }}>
        <td>{{ .Name }}</td>
        <td><select name="Attack" id="Attack"></select></td>
        <td><select name="Peek" id="Peek"></select></td>
        </tr>
        </form>
        {{ end }}

        <script>
            var allPlayers = {}
            {{ range .PlayersByName }}   
            allPlayers["{{ .Name }}"] = "{{ .Actions }}"
            {{ end }}
            
            {{ range .PlayersByName }}        
            var form = document.getElementById("{{ .Name }}")
            var attackSelect = form.getElementById("Attack")
            var peekSelect = form.getElementById("Peek")
            for (player in allPlayers) {
                if (!allPlayers[player].includes("|K|")) {
                    if (!allPlayers[{{ .Name }}].includes("|A:"+player+"|")) {
                        var option = document.createElement("option");
                        option.text = player;
                        attackSelect.add(option);
                    }
                    if (!allPlayers[{{ .Name }}].includes("|P:"+player+"|")) {
                        var option = document.createElement("option");
                        option.text = player;
                        peekSelect.add(option);
                    }
                }
            }
            {{ end }}
        </script>
     </body>
</html>
{{end}}