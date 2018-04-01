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
        <b>Actions</b><br>
        <form name="Actions" id="Actions">
        <table>
        <tr><th>Player</th><th>Attack</th><th>Peek</th></tr>
        </table>
        </form>

        <script>
            var allPlayers = {}
            {{ range .PlayersByName }}   
            allPlayers["{{ .Name }}"] = "{{ .Actions }}"
            {{ end }}

            actionsTable = document.getElementById("Actions")
            for (performingPlayer in allPlayers) {
                row = document.createElement("tr")
                playerName = document.CreateElement("td")
                playerName.innerHTML = performingPlayer
                attackCell = document.CreateElement("td")
                attackSelect = document.CreateElement("select")
                attackSelect.name = performingPlayer + "Attack"
                peekCell = document.CreateElement("td")
                peekSelect = document.CreateElement("select")
                peekSelect.name = performingPlayer + "Peek"
                for (targetPlayer in allPlayers) {
                    if (performingPlayer != targetPlayer && !allPlayers[targetPlayer].includes("|K|")) {
                        if (!allPlayers[performingPlayer].includes("|A:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.text = targetPlayer;
                            attackSelect.add(option)
                        }
                        if (!allPlayers[performingPlayer].includes("|P:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.text = targetPlayer;
                            peek.add(option)
                        }
                    }
                }
                attackCell.innerHTML = attackSelect
                row.appendChild(playerName)
                row.appendChild(attackCell)
                row.appendChild(peekCell)
                actionsTable.appendChild(row)
            }
        </script>
     </body>
</html>
{{end}}