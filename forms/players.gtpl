<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Setup New Game</title>
    </head>
    <body>
        <form action="/players" method="post">
            <table>
            <tr><th>Add player:</th><td><input type="text" name="playerName"></td></tr>
            </table>
            <input type="submit" value="Submit">
        </form>

        <form action="/roles" method="post">
            <table>
            <tr><th>Game Name:</th><td><input type="text" name="gameName" id="gameName"></td></tr>
            <tr><th>Total Players:</th><td><span id="totalPlayers"></span></td></tr>
            <tr><th>Number of Seers:</th><td><input type="text" name="seers" id="totalSeers"></td></tr>
            <tr><th>Number of Wolves:</th><td><input type="text" name="wolves" id="totalWolves"></td></tr>
            <tr><th>Remaining Villagers:</th><td><span id="totalVillagers"></span></td></tr>
            </table>
            <input type="submit" value="Submit">
        </form>

        <form action="/start" method="post">
            <input type="submit" value="Start Game">
        </form>
 
        <table id="players">
            <tr><th>Name</th></tr>
        </table>
        <script>
            playerTable = document.getElementById("players")
            numPlayersField = document.getElementById("totalPlayers")
            numVillagersField = document.getElementById("totalVillagers")
            numSeersField = document.getElementById("totalSeers")
            numWolvesField = document.getElementById("totalWolves")

            fetch("/roles")
            .then(response => response.json())
            .then(rolesList => {
                numPlayersField.innerHTML = rolesList.totalPlayers
                numVillagersField.innerHTML = rolesList.totalVillagers
                numSeersField.value = rolesList.totalSeers
                numWolvesField.value = rolesList.totalWolves
            })

            fetch("/players")
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
        </script>
    </body>
</html>