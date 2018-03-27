{{define "gameSetup.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Setup New Game</title>
    </head>
    <body>
        <form name="gameForm" onsubmit="return validateGameForm()" action="/setupGame" method="post" autocomplete="off">
            <table>
            <tr><th>Game Name:</th><td><input onkeyup="validateGameName()" type="text" name="gameName" id="gameName" autofocus></td></tr>
            <tr><td></td><td id="gameNameAlert"></td><td></td></tr>
            {{ range $key, $value := .Roles }}
            <tr><th>{{ $key }}:</th><td><input onkeyup="validateInput('{{ $key }}')" type="text" name="{{ $key }}" id="{{ $key }}" value="{{ $value }}"></td></tr>
            <tr><td></td><td id="{{ $key }}Alert"></td><td></td></tr>
            {{ end }}
            <tr><th>Total Players:</th><td><span id="totalPlayers"></span></td></tr>
            <tr><th>Symmetry (0-100)%:</th><td><input onkeyup="validateKeep()" type="text" name="keep" id="keepPercent"></td></tr>
            <tr><td></td><td id="keepAlert"></td><td></td></tr>
            <tr><td></td><td><input type="submit" value="Start Game"></td></tr>
            </table>
        </form>

        <table id="players">
            <tr><th>Name</th></tr>
        </table>

        <script>
            function TryParseInt(str,defaultValue) {
                var retValue = defaultValue;
                if(str !== null) {
                    if(str.length > 0) {
                        if (!isNaN(str)) {
                            retValue = parseInt(str);
                        }
                    }
                }
                return retValue;
            }
            function formatValidation(test, alertField, inputField, alertMessage)
            {
                if (test) {
                    document.getElementById(alertField).innerHTML = ""
                    document.getElementById(alertField).style.color = "black"
                    document.getElementById(inputField).style.borderColor = "black"
                    return true
                } else {
                    document.getElementById(alertField).innerHTML = alertMessage
                    document.getElementById(alertField).style.color = "red"
                    document.getElementById(inputField).style.borderColor = "red"
                    document.getElementById(inputField).focus()
                    return false
                }
            }

            playerTable = document.getElementById("players")
            numPlayersField = document.getElementById("totalPlayers")
            keepField = document.getElementById("keepPercent")
            gameName = document.forms["gameForm"]["gameName"].value
            numPlayers = 0
            keepPercent = 0

            fetch("/setupGame")
            .then(response => response.json())
            .then(rolesList => {
                numPlayersField.innerHTML = rolesList.totalPlayers
                numPlayers = rolesList.totalPlayers
                keepField.value = rolesList.keepPercent

                keepPercent = rolesList.keepPercent

                {{ range $key, $value := .Roles }}
                document.getElementById({{ $key }}).innerHTML = {{ $value }}
                {{ end }}
            })

            fetch("/setupPlayers")
            .then(response => response.json())
            .then(playersList => {
                //Once we fetch the list, we iterate over it
                playersList.forEach(player => {
                // Create the table row
                row = document.createElement("tr")

                // Create the table data elements for the species and
                // description columns
                playerName = document.createElement("td")
                playerName.innerHTML = player.playerName

                // Add the data elements to the row
                row.appendChild(playerName)
                // Finally, add the row element to the table itself
                playerTable.appendChild(row)
                })
            })

            function updateRoleAmounts() {
                specialRoleAmount = 0
                defaultRoleName = {{ .DefaultRoleName }}
                {{ range $key, $value := .Roles }}
                if ("{{ $key }}" != defaultRoleName) {
                    specialRoleAmount += TryParseInt(document.getElementById({{ $key }}).value, 0)
                }
                {{ end }}
                remainingDefaultRole = numPlayers - specialRoleAmount
                if (remainingDefaultRole >= 0) {
                    document.getElementById(defaultRoleName).value = (numPlayers - specialRoleAmount)
                }
            }

            function validateInput(inputType) {
                updateRoleAmounts()
                test = true

                if (test) {
                    test = numPlayers >= TryParseInt(document.getElementById(inputType).value, 0)
                    formatValidation(test, inputType + "Alert", inputType, "Must not have more " + inputType + " than total players")
                }

                if (test) {
                    test = !isNaN(document.getElementById(inputType).value) && TryParseInt(document.getElementById(inputType).value, 0) >= 0
                    formatValidation(test, inputType + "Alert", inputType, "Insert a valid number")
                }

                return test
            }

            function validateKeep() {
                keepPercent = TryParseInt(keepField.value, 0)
                return formatValidation((!isNaN(keepField.value) && keepPercent >= 0 && 100 >= keepPercent), "keepAlert", "keep", "Symmetry percent must be between 0 and 100 inclusive")
            }

            function validateRoleTotals() {
                updateRoleAmounts()
                roleTotal = 0
                {{ range $key, $value := .Roles }}
                roleTotal += TryParseInt(document.getElementById({{ $key }}).value, 0)
                {{ end }}

                test = (numPlayers == roleTotal)

                {{ range $key, $value := .Roles }}
                formatValidation(test, "{{ $key }}" + "Alert", "{{ $key }}", "Must not have more roles than total players")
                {{ end }}
                return test
            }

            function validateGameName() {
                gameName = document.forms["gameForm"]["gameName"].value
                return formatValidation((gameName != ""), "gameNameAlert", "gameName", "Must input a game name")
            }

            function validateGameForm() {
                return (validateGameName() && validateRoleTotals())
            }
        </script>
    </body>
</html>
{{end}}