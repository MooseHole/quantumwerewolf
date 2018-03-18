<html>
    <head>
    <title>Add Players</title>
    </head>
    <body>
        <form action="/player" method="post">
            Add player:
            <input type="text" name="player">
            <input type="submit" value="Submit">
            <P>
            Number of Players: <input type="text" name="total" readonly><br>
            Number of Seers: <input type="text" name="seers"><br>
            Number of Wolves: <input type="text" name="wolves"><br>
        </form>

        <table id="players">
            <tr><th>Number</th><th>Name</th></tr>
        </table>
        <script>
            playerTable = document.getElementById("players")
            //playerTable = document.querySelector("table")
        
            fetch("/player")
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
                name = document.createElement("td")
                name.innerHTML = player.name

                // Add the data elements to the row
                row.appendChild(number)
                row.appendChild(name)
                // Finally, add the row element to the table itself
                playerTable.appendChild(row)
                })
            })
        </script>
    </body>
</html>