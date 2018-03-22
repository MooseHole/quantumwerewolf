<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Setup New Game</title>
    </head>
    <body>
        <form name="gameForm" onsubmit="return validateGameForm()" action="/setupGame" method="post">
            <table>
            <tr><th>Game Name:</th><td><input onkeyup="validateGameName()" type="text" name="gameName" id="gameName"></td></tr>
            <tr><td></td><td id="gameNameAlert"></td><td></td></tr>
            <tr><th>Total Players:</th><td><span id="totalPlayers"></span></td></tr>
            <tr><th>Number of Seers:</th><td><input onkeyup="validateSeers()" type="text" name="seers" id="totalSeers"></td></tr>
            <tr><td></td><td id="totalSeersAlert"></td><td></td></tr>
            <tr><th>Number of Wolves:</th><td><input onkeyup="validateWolves()" type="text" name="wolves" id="totalWolves"></td></tr>
            <tr><td></td><td id="totalWolvesAlert"></td><td></td></tr>
            <tr><th>Remaining Villagers:</th><td><span id="totalVillagers"></span></td></tr>
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
                    return false
                }
            }

            playerTable = document.getElementById("players")
            numPlayersField = document.getElementById("totalPlayers")
            numVillagersField = document.getElementById("totalVillagers")
            numSeersField = document.getElementById("totalSeers")
            numWolvesField = document.getElementById("totalWolves")
            gameName = document.forms["gameForm"]["gameName"].value
            numPlayers = 0
            numSeers = 0
            numWolves = 0

            fetch("/setupGame")
            .then(response => response.json())
            .then(rolesList => {
                numPlayersField.innerHTML = rolesList.totalPlayers
                numVillagersField.innerHTML = rolesList.totalVillagers
                numSeersField.value = rolesList.totalSeers
                numWolvesField.value = rolesList.totalWolves

                numPlayers = rolesList.totalPlayers
                numSeers = rolesList.totalSeers
                numWolves = rolesList.totalWolves
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

            function validateSeers() {
                numSeers = TryParseInt(numSeersField.value, 0)
                return formatValidation((numSeers <= numPlayers), "totalSeersAlert", "totalSeers", "Must not have more seers than total players")
            }

            function validateWolves() {
                numWolves = TryParseInt(numWolvesField.value, 0)
                return formatValidation((numWolves <= numPlayers), "totalWolvesAlert", "totalWolves", "Must not have more wolves than total players")
            }

            function validateSpecialRoles() {
                numSeers = TryParseInt(numSeersField.value, 0)
                numWolves = TryParseInt(numWolvesField.value, 0)
                test = ((numSeers + numWolves) <= numPlayers)

                formatValidation(test, "totalSeersAlert", "totalSeers", "Must not have more wolves and seers than total players")
                formatValidation(test, "totalWolvesAlert", "totalWolves", "Must not have more wolves and seers than total players")
                return test
            }

            function validateGameName() {
                gameName = document.forms["gameForm"]["gameName"].value
                return formatValidation((gameName != ""), "gameNameAlert", "gameName", "Must input a game name")
            }

            function validateGameForm() {
                return (validateGameName() && validateSpecialRoles())
            }
        </script>
    </body>
</html>