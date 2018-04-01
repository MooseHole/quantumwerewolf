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
        <td><select name="Attack"></select></td>
        <td><select name="Peek"></select></td>
        </tr>
        </form>
        {{ end }}

        <script>
            var allPlayers = {}
            {{ range .PlayersByName }}   
            allPlayers["{{ .Name }}"] = "{{ .Actions }}"
            {{ end }}
            
            for (performingPlayer in allPlayers) {
                for (targetPlayer in allPlayers) {
                    if (performingPlayer != targetPlayer && !allPlayers[targetPlayer].includes("|K|")) {
                        var form = document.getElementById(performingPlayer)
                        var attackSelect = form.elements["Attack"]
                        if (!allPlayers[performingPlayer].includes("|A:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.text = targetPlayer;
                            attackSelect.add(option)
                        }
                        var peekSelect = form.elements["Peek"]
                        if (!allPlayers[performingPlayer].includes("|P:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.text = targetPlayer;
                            peekSelect.add(option)
                        }
                    }
                }
            }
        </script>
     </body>
</html>
{{end}}