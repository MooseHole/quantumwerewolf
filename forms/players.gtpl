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

            function validatePlayerForm() {
                var playerName = document.forms["playerForm"]["playerName"].value;
                if (playerName.length != 0 && !(/^[a-z0-9]+$/i.test(playerName))) {
                    document.getElementById("playerNameAlert").innerHTML = "Must input a valid name"
                    document.getElementById("playerNameAlert").style.color = "red"
                    document.getElementById("playerNameField").style.borderColor = "red"
                    document.getElementById("playerNameField").focus()
                    return false;
                }
            }

            function validateStartGameSettingsForm() {
                // Subtract 1 for the heading row
                var totalPlayers = document.getElementById("players").rows.length - 1

                if (totalPlayers < minPlayers) {
                    document.getElementById("gameSettingsAlert").innerHTML = "Only " + totalPlayers + " players defined.  Need at least " + minPlayers
                    document.getElementById("gameSettingsAlert").style.color = "red"
                    document.getElementById("playerNameField").style.borderColor = "red"
                    document.getElementById("playerNameField").focus()
                    return false;
                }
            }
        </script>
    </body>
</html>