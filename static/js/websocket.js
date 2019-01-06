var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws/join');
    // Message received on the socket
    socket.onmessage = function (event) {

        var data = JSON.parse(event.data);
        
        console.log(data);

        switch (data.Type) {
        case 0: // JOIN
            break;
        case 1: // LEAVE
            break;
        case 2: // MESSAGE
            d = new Date();
            hrs = d.getHours().toString();
            mins = d.getMinutes().toString();
            secs = d.getSeconds().toString();
            time = hrs + ":" + mins + ":" + secs;

            chart.data.labels.push(time);
            chart.data.datasets[0].data.push(data.Content);
            if (parseInt(data.Content) > chart.options.scales.yAxes[0].ticks.max) {
                chart.options.scales.yAxes[0].ticks.max = parseInt(data.Content);
            }
            chart.update();

            $('#stats').text(data.Content + " ms");
            break;
        case 3: // TOTAL REQUESTS
            $('#requests').text(data.Content);
            break;
        case 4: // P90
            $('#ninety').text(data.Content + " ms");
            break;
        case 5: //P99
            $('#ninetynine').text(data.Content + " ms");
            break;
        case 6: //P50
            $('#fifty').text(data.Content + " ms");
            break;
        }
    };

    $("#stop").click(function () {
        url = window.location.protocol + '//' + window.location.host + '/?command=stop'
        $.ajax({url: url, type: "POST", success: function(result){
            $("#toast").show();;
        }});
    });

    var ctx = document.getElementById("myAreaChart");
    var chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
            label: "Sessions",
            lineTension: 0.3,
            backgroundColor: "rgba(2,117,216,0.2)",
            borderColor: "rgba(2,117,216,1)",
            pointRadius: 5,
            pointBackgroundColor: "rgba(2,117,216,1)",
            pointBorderColor: "rgba(255,255,255,0.8)",
            pointHoverRadius: 5,
            pointHoverBackgroundColor: "rgba(2,117,216,1)",
            pointHitRadius: 50,
            pointBorderWidth: 2,
            data: [],
            }],
        },
        options: {
            scales: {
                xAxes: [{
                    time: {
                    unit: 'date'
                    },
                    gridLines: {
                    display: false
                    },
                    ticks: {
                    maxTicksLimit: 7
                    }
                }],
                yAxes: [{
                    ticks: {
                        min: 0,
                        max: 50,
                        maxTicksLimit: 5
                    },
                    gridLines: {
                        color: "rgba(0, 0, 0, .125)",
                    }
                }],
            },
            legend: {
            display: false
            }
        }
    });
});

Chart.defaults.global.defaultFontFamily = '-apple-system,system-ui,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif';
Chart.defaults.global.defaultFontColor = '#292b2c';

