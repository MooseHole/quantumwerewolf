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
        </table>
        </form>

        <ul>
        {{ range .ActionSubjects }}
            <li>{{ .Subject }}
                <ul>
                {{ range .Peek }}
                    <li>{{ . }}</li>
                {{ end }}
                </ul>
                <ul>
                {{ range .Attack }}
                    <li>{{ . }}</li>
                {{ end }}
                </ul>
                <ul>
                {{ range .Lynch }}
                    <li>{{ . }}</li>
                {{ end }}
                </ul>
            </li>
        {{ end }}
        </ul>
        

        <script>
            var allPlayers = {}
            {{ range .PlayersByName }}   
            allPlayers["{{ .Name }}"] = "{{ .Actions }}"
            {{ end }}

            actionsForm = document.getElementById("Actions")
            actionsTable = document.getElementById("ActionsTable")
            for (performingPlayer in allPlayers) {
                if (allPlayers[performingPlayer].includes("#")) {
                    continue
                }
                row = document.createElement("tr")
                playerName = document.createElement("td")
                playerName.innerHTML = performingPlayer
                row.appendChild(playerName)
                {{ if .IsNight }}
                attackCell = document.createElement("td")
                attackSelect = document.createElement("select")
                attackSelect.id = performingPlayer + "Attack"
                attackSelect.name = performingPlayer + "Attack"
                peekCell = document.createElement("td")
                peekSelect = document.createElement("select")
                peekSelect.name = performingPlayer + "Peek"
                var noAttackOption = document.createElement("option");
                noAttackOption.value = "";
                noAttackOption.text = "--NONE--";
                attackSelect.appendChild(noAttackOption)
                var noPeekOption = document.createElement("option");
                noPeekOption.value = "";
                noPeekOption.text = "--NONE--";
                peekSelect.appendChild(noPeekOption)
                for (targetPlayer in allPlayers) {
                    if (performingPlayer != targetPlayer && !allPlayers[targetPlayer].includes("#")) {
                        if (!allPlayers[performingPlayer].includes("@"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.value = targetPlayer;
                            option.text = targetPlayer;
                            attackSelect.appendChild(option)
                        }
                        if (!allPlayers[performingPlayer].includes("%"+targetPlayer+"|")) {
                            var option = document.createElement("option");
                            option.value = targetPlayer;
                            option.text = targetPlayer;
                            peekSelect.appendChild(option)
                        }
                    }
                }
                attackSelect.form = "Actions"
                peekSelect.form = "Actions"
                attackCell.appendChild(attackSelect)
                peekCell.appendChild(peekSelect)
                row.appendChild(attackCell)
                row.appendChild(peekCell)
                {{ else }}
                lynchCell = document.createElement("td")
                lynchSelect = document.createElement("select")
                lynchSelect.id = performingPlayer + "Lynch"
                lynchSelect.name = performingPlayer + "Lynch"
                var noLynchOption = document.createElement("option");
                noLynchOption.value = "";
                noLynchOption.text = "--NONE--";
                lynchSelect.appendChild(noLynchOption)
                for (targetPlayer in allPlayers) {
                    if (performingPlayer != targetPlayer && !allPlayers[targetPlayer].includes("#")) {
                        var option = document.createElement("option");
                        option.value = targetPlayer;
                        option.text = targetPlayer;
                        lynchSelect.appendChild(option)
                    }
                }
                lynchSelect.form = "Actions"
                lynchCell.appendChild(lynchSelect)
                row.appendChild(lynchCell)
                {{ end }}
                actionsTable.appendChild(row)
            }
        </script>
     </body>
</html>
{{end}}