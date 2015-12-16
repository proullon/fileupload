<html>
<head>
    <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
      <input type="file" name="uploadfile" />
      <input type="submit" value="upload" />
</form>


{{ range . }}
<p>
<img src="/images/{{.}}">
</p>
{{end}}

</body>
</html>

