<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Meteor Burst</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        <link rel="stylesheet" href="/static/css/bootstrap/bootstrap.min.css">
        <link rel="stylesheet" href="/static/fontawesome-free/css/all.min.css">
        <link rel="stylesheet" href="/static/datatables/dataTables.bootstrap4.css">
        <link rel="stylesheet" href="/static/css/sb-admin.css">
        {{ block "css" . }}{{ end }}
    </head>
    <body id="page-top">
        <nav class="navbar navbar-expand navbar-dark bg-dark static-top">
            <a class="navbar-brand mr-1" href="/">Meteor Burst</a>
        </nav>
        <div id="content-wrapper" class="container" style="position:relative;top:30px;">
            {{ block "content" . }}{{ end }}
        </div>
        <script src="static/js/jquery/jquery.min.js"></script>
        <script src="/static/js/bootstrap/bootstrap.min.js"></script>
        <script src="/static/chart.js/Chart.min.js"></script>
        {{ block "js" . }}{{ end }}
    </body>
</html>