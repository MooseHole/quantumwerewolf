<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        I'm a game!
        <script>
            fetch("/game")
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