{{define "game.gtpl"}}<!DOCTYPE html>
<html lang="en">
    <head>
    <title>Game</title>
    </head>
    <body>
        <img src="data:image/png;base64,{{ .Graph }}" width="200" height="200">
        Game: {{ .Name }}<br>
        {{ .WinMessage }}<br>
        Remaining Universes: {{ len .Universes }}<br>
        Number of Players: {{ .TotalPlayers }}<br>
        {{ range $name, $value := .Roles }}
        Number of {{ $name }}: {{ $value }}<br>
        {{ end }}
        Round: {{ .Round }}
        <p>
        [table][tr][td][b]Player[/b][/td][td][b]Good[/b][/td][td][b]Evil[/b][/td][td][b]Dead[/b][/td][td][b]Name[/b][/td][td][b]Role[/b][/td][/tr]
        {{ range $num, $selections := .ActionSubjects }}
            {{ $thisGood := index $selections.Percents "Good" }}
            {{ $thisEvil := index $selections.Percents "Evil" }}
            {{ $thisDead := index $selections.Percents "Dead" }}
            [tr][td]{{ $num }}[/td][td]{{ $thisGood }}%[/td][td]{{ $thisEvil }}%[/td][td]{{ $thisDead }}%[/td][td]{{ $selections.RevealName }}[/td][td]{{ $selections.RevealRole }}[/td][/tr]
        {{ end }}
        [/table]
        <p>
        <table>
        <tr><th>Player</th><th>Good</th><th>Evil</th><th>Dead</th><th>Name</th><th>Role</th></tr>
        {{ range $num, $selections := .ActionSubjects }}
            {{ $thisGood := index $selections.Percents "Good" }}
            {{ $thisEvil := index $selections.Percents "Evil" }}
            {{ $thisDead := index $selections.Percents "Dead" }}
            <tr><td>{{ $num }}</td><td>{{ $thisGood }}%</td><td>{{ $thisEvil }}%</td><td>{{ $thisDead }}%</td><td>{{ $selections.RevealName }}</td><td>{{ $selections.RevealRole }}</td></tr>
        {{ end }}
        </table>
        <p>
        <table>
        {{ range $num, $selections := .ActionSubjects }}
            {{ if eq $num 0 }}
                <tr><th>Player</th><th>Name</th>
                {{ range $percentName, $percentSelection := $selections.Percents }}
                    <th>{{ $percentName }}</th>
                {{ end }}
                </tr>
            {{ end }}
        {{ end }}
        {{ range $num, $selections := .ActionSubjects }}
            <tr><td>{{ $num }}</td><td>{{ $selections.Name }}</td>
            {{ range $percentName, $percentSelection := $selections.Percents }}
                <td>{{ $percentSelection }}%</td>
            {{ end }}
            </tr>
        {{ end }}
        </tr>
        </table>
        <p>
        <b>Action Messages</b><br>
        {{ range $num, $messages := .ActionMessages }}
            {{ range $messages }}
                {{ . }}<br>
            {{ end }}
            ------------------<br>
        {{ end }}
        <p>
        <b>Actions for {{ .Round }}</b><br>
        <form name="Actions" id="Actions" action="/processActions" method="post">
        <input type="hidden" name="gameId" value={{ .GameID }}>
        <input type="submit">
        <input type="checkbox" name="advance" value="true">Advance to next round<br>
        <table border=1 name="ActionsTable" id="ActionsTable">
        <tr><th></th><th></th>
        {{ range .Rounds }}
        <th colspan=4>Round {{ . }}</th>
        {{ end }}
        </tr>
        <tr><th>Player</th><th>Name</th>
        {{ range .Rounds }}
        <th>Vote</th><th>Attack</th><th>Peek</th><th>Died</th>
        {{ end }}
        </tr>
        {{ range $num, $selections := .ActionSubjects }}
        <tr>
            <td>{{ $num }}</td><td>{{ $selections.Name }}</td>
            {{ range $roundNum := $.Rounds }}
                {{ $thisAttacked := index $selections.Attacked $roundNum }}
                {{ $thisPeeked := index $selections.Peeked $roundNum }}
                {{ $thisPeekResult := index $selections.PeekResult $roundNum }}
                {{ $thisVoted := index $selections.Voted $roundNum }}
                {{ $thisKilled := index $selections.Killed $roundNum }}
                
                <td>
                {{ if and (not (eq (len $selections.Vote) 1)) (and (eq $.RoundNum $roundNum) (not $.IsNight)) }}
                    <select name="{{ $num }}Vote">
                    {{ range $name, $value :=  $selections.Vote }}
                        <option value="{{ $value }}"
                        {{ $thisSelection := index $selections.Voted $roundNum }}
                        {{ if eq $value $thisSelection }}
                        selected
                        {{ end }}
                        >{{ $name }}</option>
                    {{ end }}
                    </select>
                {{ else }}
                    {{ $thisVoted }}
                {{ end }}
                </td>
                {{ if and (not (eq (len $selections.Attack) 1)) (and (eq $.RoundNum $roundNum) $.IsNight) }}
                <td>
                    <select name="{{ $num }}Attack">
                    {{ range $name, $value :=  $selections.Attack }}
                        <option value="{{ $value }}"
                        {{ $thisSelection := index $selections.Attacked $roundNum }}
                        {{ if eq $value $thisSelection }}
                        selected
                        {{ end }}
                        >{{ $name }}</option>
                    {{ end }}
                    </select>
                </td>
                {{ else }}
                <td>{{ $thisAttacked }}</td>
                {{ end }}
                {{ if and (not (eq (len $selections.Peek) 1)) (and (eq $.RoundNum $roundNum) $.IsNight) }}
                <td>
                    <select name="{{ $num }}Peek">
                    {{ range $name, $value :=  $selections.Peek }}
                        <option value="{{ $value }}"
                        {{ $thisSelection := index $selections.Peeked $roundNum }}
                        {{ if eq $value $thisSelection }}
                        selected
                        {{ end }}
                        >{{ $name }}</option>
                    {{ end }}
                    </select>
                </td>
                {{ else }}
                <td>{{ $thisPeeked }}{{ $thisPeekResult }}</td>
                {{ end }}
                <td>{{ $thisKilled }}</td>
            {{ end }}
        </tr>
        {{ end }}
        </table>
        </form>

        <b>Actions</b><br>
        {{ range .PlayersByNum }}
            Player {{ .Num }} ({{ .Name }}): {{ .Actions }}<br>
        {{ end }}
        <p>
        <b>Multiverse</b><br>
        {{ range $num, $output := .Universes }}
        {{ $output  }} {{ $num }}<br>
        {{ end }}
     </body>
</html>
{{end}}