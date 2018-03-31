{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        Game: {{ .Name }}<br>
        Number of Players: {{ .TotalPlayers }}<br>
        Round: {{ .IsNight }} {{ .Round }}<p>
        <b>Actions</b><br>
        {{ range .PlayersByNum }}
            Player {{ .Num }}: {{ .Actions }}<br>
        {{ end }}    
        <p>
        <b>Players</b><br>
        <table>
        {{ range .PlayersByName }}        
            <tr>
            <form name={{ .Name }}>
            <td>{{ .Name }}</td>
            <td><select name="Attack"></select></td>
            <td><select name="Peek"></select></td>
            </tr>
            </form>
        {{ end }}

        <script>
            var selects = document.getElementsByTagName('select');
            for(var z=0; z<selects.length; z++){
                {{ range .PlayersByName }}        
                {
                    var option = document.createElement("option");
                    option.text = "{{ .Name }}";
                    selects[z].add(option);
                }
                {{ end }}
            }
        </script>
     </body>
</html>
{{end}}