package web

import (
	"html/template"
)

var Template *template.Template

func init() {
	Template = template.Must(template.New("").Parse(indexTemplate))
}

const indexTemplate = `<html>
<head>
  <title>Build dashboard</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.css">
  <script src="https://code.jquery.com/jquery-3.4.1.min.js" integrity="sha256-CSXorXvZcTkaix6Yvo6HppcZGetbYMGWSFlBw8HfCJo=" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.js"></script>
  <style>
    i.amber.icon {color: #FFA000;}
  </style>
</head>
<body>

<div class="ui menu inverted huge">
  <div class="header item">
    Dashboard
  </div>
  <div class="ui item dropdown simple">
    Repositories<i class="dropdown icon"></i>
    <div class="menu">
      {{- range .Repositories }}
      <a class="item">{{ .Name }}</a>
      {{- end }}
      <div class="ui divider"></div>
      <a class="item" onclick="newRepository();">New...</a>
      <div class="ui item dropdown simple">
        <i class="tiny trash icon"></i>Delete
        <div class="menu">
        {{- range .Repositories }}
        <a class="item" onclick="deleteRepository({{ .Id }}, '{{ .Name }}');">{{ .Name }}</a>
        {{- end }}
        </div>
      </div>
    </div>
  </div>
  <div class="ui item dropdown simple">
    Job<i class="dropdown icon"></i>
    <div class="menu">
      {{- range .Jobs }}
      <a class="item">{{ .Command }} {{ .Repository.Name }}{{ .Target }}</a>
      {{- end }}
    </div>
  </div>
</div>

<!-- modal -->
<div class="ui newRepo modal">
  <i class="close icon"></i>
  <div class="header">
    New Repository
  </div>
  <div class="content">
    <form class="ui form newRepo">
      <div class="field">
        <label>Name</label>
        <input type="text" name="name" placeholder="The name of the repository">
      </div>
      <div class="field">
        <label>URL</label>
        <input type="text" name="url" placeholder="URL of the repository (e.g https://github.com/f110/sandbox)">
      </div>
      <div class="field">
        <label>Clone URL</label>
        <input type="text" name="clone_url" placeholder="URL for cloning of the repository (e.g https://github.com/f110/sandbox.git)">
      </div>
      <button class="ui button" onclick="createRepository()">Add</button>
    </form>
  </div>
</div>

<div class="ui basic modal">
  <div class="ui icon header">
    <i class="archive icon"></i>
    Delete "<span class="repoName"></span>" repository
  </div>
  <div class="actions">
    <div class="ui red basic cancel inverted button">
      <i class="remove icon"></i>
      No
    </div>
    <div class="ui green ok inverted button">
      <i class="checkmark icon"></i>
      Yes
    </div>
  </div>
</div>
<!-- end of modal -->

<div class="ui container">
  {{- range .Jobs }}
  <h3 class="ui block header">
    <div class="ui grid">
      <div class="two column row">
        <div class="left floated column">{{ if .Success }}<i class="green check icon"></i>{{ else }}<i class="red attention icon"></i>{{ end }}{{ .Command }} {{ .Repository.Name }}{{ .Target }}</div>
        <div class="right aligned floated column"><a href="{{ $.APIHost }}/run?job_id={{ .Id }}"><i class="green play icon"></i></a></div>
      </div>
    </div>
  </h3>

  <div class="ui container">
    <table class="ui selectable striped table">
      <thead>
        <tr>
          <th>#</th>
          <th>OK</th>
          <th>Log</th>
          <th>Rev</th>
          <th>Trigger</th>
          <th>Start at</th>
          <th>Finished at</th>
        </tr>
      </thead>
      <tbody>
		{{- range .Tasks }}
        <tr>
          <td>{{ .Id }}</td>
          <td>{{ if .FinishedAt }}{{ if .Success }}<i class="green check icon"></i>{{ else }}<i class="red attention icon"></i>{{ end }}{{ else }}<i class="sync amber alternate icon"></i>{{ end }}</td>
          <td>{{ if .LogFile }}<a href="/logs/{{ .LogFile }}">text</a>{{ end }}</td>
          <td><a href="">{{ .Revision }}</a></td>
          <td>{{ .Via }}</td>
          <td>{{ .CreatedAt.Format "2006/01/02 15:04:06" }}</td>
          <td>{{ if .FinishedAt }}{{ .FinishedAt.Format "2006/01/02 15:04:06" }}{{ end }}</td>
        </tr>
        {{- end }}
      </tbody>
    </table>
    {{- end }}
  </div>
</div>

<script>
function newRepository() {
	$('.ui.newRepo.modal').modal({centered:false}).modal('show');
}

function createRepository() {
	var f = document.querySelector('.ui.form.newRepo');
	var params = new URLSearchParams();
	params.append("name", f.name.value);
	params.append("url", f.url.value);
	params.append("clone_url", f.clone_url.value);
	fetch('/new_repo', {
		method: 'POST',
		body: params,
	});
}

function deleteRepository(id, name) {
  var e = document.querySelector('span.repoName');
  e.textContent = name;
  $('.ui.basic.modal').modal({
    onApprove: function() {
		var params = new URLSearchParams();
		params.append("id", id);
		fetch('/delete_repo', {
			method: 'POST',
			body: params,
		});
	},
  }).modal('show');
}
</script>
</body>
</html>`