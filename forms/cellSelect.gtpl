{{ define "cellSelect" }}
<td>
<select id={{ .ID }} name={{ .Name }}>
    {{ range .Options }}
    <option value={{ .Value }}>{{.Text}}</option>
    {{ end }}
</select>
</td>
{{ end }}
