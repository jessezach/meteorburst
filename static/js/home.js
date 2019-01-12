var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws/join');
    // Message received on the socket
    socket.onmessage = function (event) {

        var data = JSON.parse(event.data);
        
        console.log(data);

        switch (data.Type) {
        case 8: //SLAVES
            $('#slaves').text(data.Content);
            break;
        }
    };
});