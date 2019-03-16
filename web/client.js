const {MessageRequest, MessageResponse, AuthenticateRequest, AuthenticateResponse} = require('./flipapps_pb.js');
const {GetInfoRequest, GetInfoResponse} = require('./flipdot_pb.js');
const {FlipAppsClient} = require('./flipapps_grpc_web_pb.js');

var client = new FlipAppsClient('https://jimsflipdot.hopto.org');
var request = new GetInfoRequest();

console.log("Initialising...")
var authForm = document.getElementById("auth-form");
var messageForm = document.getElementById("message-form");

var token;

authForm.addEventListener('submit', function(event) {
    console.log("Password submitted!");
    event.preventDefault(); // Disable refresh of page
    // Get submitted form
    data = new FormData(event.target);
    console.log(data);
    // Construct a request
    var request = new AuthenticateRequest();
    request.setPassword(data.get('password'));
    // Send the request
    client.authenticate(request, {}, (err, response) => {
        // Save the token globally
        token = response.token;
        console.log(token);
    });
});

messageForm.addEventListener('submit', function(event) {
    console.log("submitted!");
    // Disable refresh of page
    event.preventDefault();
    // Get submitted text
    data = new FormData(event.target);
    console.log(data);
    // Construct a request
    var request = new MessageRequest();
    request.setFrom(data.get('from'));
    request.setText(data.get('text'));
    // Send the request
    client.sendMessage(request, {'token': token}, (err, response) => {
        console.log(response);
    });
});
