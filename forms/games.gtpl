<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Games</title>
    </head>
    <body>
        <form action="/setupPlayers" method="post">
            <input type="submit" value="New Game">
        </form>
        <table id="games">
            <tr><th>Name</th></tr>
        </table>
        <script>
            gameTable = document.getElementById("games")

            fetch("/getGames")
            .then(response => response.json())
            .then(gameList => {
                gameList.forEach(game => {
                    row = document.createElement("tr")

                    gameName = document.createElement("td")
                    gameName.innerHTML = game.gameName

                    row.appendChild(gameName)
                    gameTable.appendChild(row)
                })
            })
        </script>
    </body>
</html>