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
        </form>
        <table name="ActionsTable" id="ActionsTable">
        <tr><th>Player</th><th>Attack</th><th>Peek</th></tr>
        </table>

        <script>
            var allPlayers = {}
            {{ range .PlayersByName }}   
            allPlayers["{{ .Name }}"] = "{{ .Actions }}"
            {{ end }}

            actionsForm = document.getElementById("Actions")
            actionsTable = document.getElementById("ActionsTable")
            for (performingPlayer in allPlayers) {
                row = document.createElement("tr")
                playerName = document.createElement("td")
                playerName.innerHTML = performingPlayer
                attackCell = document.createElement("td")
                attackSelect = document.createElement("select")
                attackSelect.id = performingPlayer + "Attack"
                attackSelect.name = performingPlayer + "Attack"
                peekCell = document.createElement("td")
                peekSelect = document.createElement("select")
                peekSelect.name = performingPlayer + "Peek"
                for (targetPlayer in allPlayers) {
                    if (performingPlayer != targetPlayer && !allPlayers[targetPlayer].includes("|K|")) {
                        if (!allPlayers[performingPlayer].includes("|A:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.value = targetPlayer;
                            option.text = targetPlayer;
                            attackSelect.appendChild(option)
                        }
                        if (!allPlayers[performingPlayer].includes("|P:"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.value = targetPlayer;
                            option.text = targetPlayer;
                            peekSelect.appendChild(option)
                        }
                    }
                }
                attackSelect.form = "Actions"
                peekSelect.form = "Actions"
                attackCell.innerHTML = attackSelect
                peekCell.innerHTML = peekSelect
                row.appendChild(playerName)
                row.appendChild(attackCell)
                row.appendChild(peekCell)
                actionsTable.appendChild(row)
            }
        </script>
     </body>
</html>
{{end}}