package main

var templateIndex = []byte(`
<html>
  <head>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <script   src="https://code.jquery.com/jquery-3.1.1.min.js"   integrity="sha256-hVVnYaiADRTO2PzUGmuLJr8BLUSjGIZsDYGmIJLv2b8="   crossorigin="anonymous"></script>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
  <script src="http://momentjs.com/downloads/moment.min.js"></script>
  <style>
    .name {
      font-weight:300;
    }
    .value {
      font-size:16px;
      text-color:#666;
      padding-left:20px;
    }
    .value.positive {
      color:green;
      font-weight:bold;
    }
    .value.positive::before { 
        content: "+";
    }
    .positive::before { 
        content: "+";
    }
    .time {
      vertical-align:middle !important;
      text-align:center;
      font-size:85%;
    }
    .raw {
      background-color: #333;
      color: #ddd;
      font-family:monospace;
      font-size:80%;
      display:none;
    }
  </style>
    <script>
      $(document).ready(function(){
        var url = "ws://" + window.location.host + "/ws";
        var ws = new WebSocket(url);
        var records = $("#records table tbody");

        function addRecord(record){
          var tr = $("<tr>").addClass("record");
          var time = $("<td>").addClass("time");
          var timeBig = $("<span>").addClass("hidden-xs");
          var timeSmall = $("<span>").addClass("visible-xs-* hidden-sm hidden-md hidden-lg");
          
          var d = moment(new Date(record.time));

          timeBig.text(d.format("MMM Do YY, hh:mm:ss"));
          timeSmall.text(d.format("hh:mm:ss"));
          time.append(timeBig,timeSmall);

          tr.append(time);
          $.each(record.locations,function(i,location){
            var td = $("<td>");
            var loc = $("<div>").addClass("loc");

            var name = $("<div class='name'>").text(location.name);
            var value = $("<div class='value'>").text(location.value);

            if(location.status=="*") {
              td.addClass("success");
            }

            if(location.positive){
              value.addClass("positive");
            } else {
              value.prepend(location.positive_value);
            }

            loc.append(name,value);
            tr.append(td.append(loc));
          });


          var rawRecord = $("<td>").text(record.raw);
          rawRecord.addClass("raw").attr("colspan",record.locations.length+1);
          records.prepend($("<tr>").append(rawRecord));

          tr.hide();
          records.prepend(tr);
          tr.fadeIn();
        }

        ws.onmessage = function (msg) {
          var record = JSON.parse(msg.data)
          addRecord(record);
        };

        $("#records").on("click","td",function(){
          var tr = $(this).closest("tr");
          tr.siblings("tr").find("td.raw").hide();
          tr.next('tr').find("td.raw").toggle();
        });
        $("#records").on("click","td.raw",function(){
          $(this).hide();
        });
      });  
    </script>
  </head>
  <body>
    <div id="records" class="container-fluid">
      <table class='table table-hover table-bordered table-condensed'><tbody></tbody></table>
    </div>
  </body>
</html>
`)
