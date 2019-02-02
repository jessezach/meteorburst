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

    $('#ramp-type').change(function(){
        $('#ramp-field').children().remove();
        if( $(this).val() == 'linear'){
            $('#ramp-field').append('<input id="ramp" name="ramp" class="form-control" placeholder="seconds" type="text" style="position:relative;top:10px;"/>');
        }else {
            $('#ramp-field').append('<textarea class="form-control" id="step" name="step" rows="6" style="position:relative;top:10px;"></textarea>')
        }
    });
});