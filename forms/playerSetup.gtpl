<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Setup New Game</title>
    </head>
    <body>
        <table>
        <form name="playerForm" action="/setupPlayers" method="post" onsubmit="return validatePlayerForm()" autocomplete="off">
        <tr><th>Add player:</th><td><input id="playerNameField" type="text" name="playerName" autofocus></td><td><input type="submit" value="Submit"></td></tr>
        <tr><td></td><td id="playerNameAlert"></td><td></td></tr>
        </form>
        <form action="/startGameSetup" method="get" onsubmit="return validateStartGameSettingsForm()">
        <tr><td></td><td><input type="submit" value="Start Game Settings"></td><td></td></tr>
        <tr><td></td><td id="gameSettingsAlert"></td><td></td></tr>
        </form>
        </table>
 
        <table id="players">
            <tr><th>Name</th></tr>
        </table>
        <script>
            var minPlayers = 3
            var maxPlayers = 20 // Because 20! is the highest factorial that fits in uint64

            playerTable = document.getElementById("players")

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

            function validatePlayerForm() {
                var playerName = document.forms["playerForm"]["playerName"].value;
                test = (playerName.length != 0 && (/^[a-z0-9]+$/i.test(playerName)))
                formatValidation(test, "playerNameAlert", "playerNameField", "Must input a valid name")
                if (test)
                {
                    // Start at 1 to skip the header
                    for (var i = 1, row; row = playerTable.rows[i]; i++) {
                        previousName = row.cells[0].innerHTML
                        if (previousName == playerName) {
                            test = false
                            break
                        }
                    }
                    formatValidation(test, "playerNameAlert", "playerNameField", "Duplicate names not allowed")
                }

                if (test) {
                    var totalPlayers = document.getElementById("players").rows.length - 1
                    test = (totalPlayers < maxPlayers)
                    formatValidation(test, "gameSettingsAlert", "playerNameField", "Only " + maxPlayers + " may be added.")
                }

                return test
            }

            function validateStartGameSettingsForm() {
                // Subtract 1 for the heading row
                var totalPlayers = document.getElementById("players").rows.length - 1
                test = (totalPlayers >= minPlayers)
                formatValidation(test, "gameSettingsAlert", "playerNameField", "Only " + totalPlayers + " players defined.  Need at least " + minPlayers)
                return test
            }
        </script>
    </body>
</html>