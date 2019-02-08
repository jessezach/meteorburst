var socket;

$(document).ready(function () {
    var i = 1;
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

    var table = '<table class="table table-sm">\
    <thead>\
      <tr>\
        <th scope="col">Users</th>\
        <th scope="col">Duration</th>\
        <th scope="col">Unit</th>\
      </tr>\
    </thead>\
    <tbody id="tbody">\
    <tr>\
      <td><input name="usr[0]" type="number" placeholder="10" required="required"></td>\
      <td><input name="dur[0]" type="number" placeholder="20" required="required"></td>\
      <td>\
        <select name="unit[0]" id="unit">\
            <option value="seconds">seconds</option>\
            <option value="minutes">minutes</option>\
        </select>\
      </td>\
    </tr>\
    </tbody>\
    </table>'

    var row = `<tr>\
    <td><input type="number" name="usr[${i}]" placeholder="10" required="required"></td>\
    <td><input type="number" name="dur[${i}]" placeholder="20" required="required"></td>\
    <td>\
      <select name="unit[${i}]" id="unit">\
        <option value="seconds">seconds</option>\
        <option value="minutes">minutes</option>\
      </select>\
    </td>\
    <tr>`;

    $('select[name="format"]').change(function() {
        $('#duration').remove();
        if( $(this).val() != 'none') {
            $('#duration-field').append('<input type="number" name="duration" class="form-control" id="duration" placeholder="20" style="position:relative;top:10px;" required="required">')
        }
    });

    $('#ramp-type').change(function() {
        $('#ramp-field').children().remove();
        $('#add-row').hide();

        if( $(this).val() == 'linear') {
            $('#ramp-field').append('<input id="ramp" name="ramp" class="form-control" placeholder="seconds" type="number" style="position:relative;top:10px;" required="required"/>');
        } else if( $(this).val() == 'step') {
            $('#ramp-field').append(table);
            $('#add-row').show();
        }
    });

    $('#add-row').click(function() {
        $('#tbody').append(row);
        i += 1;
        row = `<tr>\
            <td><input type="number" name="usr[${i}]" placeholder="10" required="required"></td>\
            <td><input type="number" name="dur[${i}]" placeholder="20" required="required"></td>\
            <td>\
            <select name="unit[${i}]" id="unit">\
                <option value="seconds">seconds</option>\
                <option value="minutes">minutes</option>\
            </select>\
            </td>\
            <tr>`;
    });
});