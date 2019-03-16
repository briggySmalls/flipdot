const {MessageRequest, MessageResponse} = require('./flipapps_pb.js');
const {GetInfoRequest, GetInfoResponse} = require('./flipdot_pb.js');
const {FlipAppsClient} = require('./flipapps_grpc_web_pb.js');

var client = new FlipAppsClient('https://jimsflipdot.hopto.org');
var request = new GetInfoRequest();

// Run 'getInfo' request immediately for fun
client.getInfo(request, {}, (err, response) => {
  console.log(response);
});

console.log("Initialising...")
var form = document.getElementById("myform");

form.addEventListener('submit', function(event) {
    console.log("submitted!");
    // Disable refresh of page
    event.preventDefault();
    // Get submitted text
    data = new FormData(form);
    console.log(data);
    // Construct a request
    var request = new MessageRequest();
    request.setFrom(data.get('from'));
    request.setText(data.get('text'));
    // Send the request
    client.sendMessage(request, {}, (err, response) => {
        console.log(response);
    });
});
