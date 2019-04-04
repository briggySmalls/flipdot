import {MessageRequest, MessageResponse, AuthenticateRequest, AuthenticateResponse} from '../generated/flipapps_pb';
import {GetInfoRequest, GetInfoResponse} from '../generated/flipdot_pb';
import {FlipAppsClient} from '../generated/FlipAppsServiceClientPb';

export class Client {
    // Token used to send authenticated gRPC messages with
    private token: string = '';

    // gRPC client
    private client: FlipAppsClient;

    constructor(domain: string) {
        // Create a flipapps client
        this.client = new FlipAppsClient(domain, null, null);
    }

    public authenticate(password: string) {
        // Construct a request
        const request = new AuthenticateRequest();
        request.setPassword(password);
        // Send the request
        this.client.authenticate(request, {}, (err, response) => {
            if (err != null) {
                return;
            }
            // Save the token globally
            this.token = response.getToken();
        });
    }

    public sendTextMessage(from: string, text: string, callback: (err: any, response: any) => void) {
        // Construct a request
        const request = new MessageRequest();
        request.setFrom(from);
        request.setText(text);
        // Send the request
        this.client.sendMessage(request, {token: this.token}, callback);
    }
}
