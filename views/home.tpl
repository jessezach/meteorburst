{{ template "base.tpl" . }}
{{ define "content" }}
    {{ if .error }}
    <div class="alert alert-danger" role="alert">
        {{ .error }}
    </div>
    {{ end }}
    <form action="/?command=start" method="post">
        <div class="form-group">
            <div class="form-row">
                <div class="col-md-6">
                    <label for="url">URL</label>
                    <input type="text" name="url" class="form-control" id="url" placeholder="https://www.example.com/users" required="required">
                </div>
            </div>
        </div>
        <div class="form-group">
            <div class="form-row">
                <div class="col-md-6">
                    <label for="headers">Headers</label>
                    <input type="text" name="headers" class="form-control" id="headers" placeholder="accept: application/json;content: application/json">
                </div>
            </div>
        </div>
        <div class="form-group">
            <div class="form-row">
                <div class="col-md-6">
                    <label for="method">Method</label>
                    <select class="form-control" name="method" id="method">
                        <option>GET</option>
                        <option>POST</option>
                        <option>PATCH</option>
                        <option>PUT</option>
                        <option>DELETE</option>
                    </select>
                </div>
            </div>    
        </div>
        <div class="form-group">
            <div class="form-row">
                <div class="col-md-6">
                    <label for="payload">Payload</label>
                    <textarea class="form-control" id="payload" name="payload" rows="6"></textarea>
                </div>
            </div>
        </div>
        <div class="form-group">
            <div class="form-row">
                <div class="col-md-6">
                    <label for="users">Users</label>
                    <input type="text" name="users" class="form-control" id="users" placeholder="50" required="required">
                </div>
            </div>
        </div>
        <button type="submit" class="btn btn-primary">Start</button>
    </form>
{{ end }}