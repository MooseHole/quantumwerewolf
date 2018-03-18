<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Add Players</title>
    </head>
    <body>
        <form action="/players" method="post">
            Add player:
            <input type="text" name="playerName">
            <input type="submit" value="Submit">
        </form>
            
        <form action="/roles" method="post">
            Number of Players: <input type="text" name="total" id="totalPlayers" readonly><br>
            Number of Villagers: <input type="text" name="villagers" id="totalVillagers" readonly><br>
            Number of Seers: <input type="text" name="seers" id="totalSeers"><br>
            Number of Wolves: <input type="text" name="wolves" id="totalWolves"><br>
            <input type="submit" value="Submit">
        </form>
 
        <table id="players">
            <tr><th>Number</th><th>Name</th></tr>
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
                numPlayersField.value = rolesList.totalPlayers
                numVillagersField.value = rolesList.totalVillagers
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
                number = document.createElement("td")
                number.innerHTML = player.number
                playerName = document.createElement("td")
                playerName.innerHTML = player.playerName

                // Add the data elements to the row
                row.appendChild(number)
                row.appendChild(playerName)
                // Finally, add the row element to the table itself
                playerTable.appendChild(row)
                })
            })
        </script>
    </body>
</html>