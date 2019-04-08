import {MessageRequest, MessageResponse, AuthenticateRequest} from '../generated/flipapps_pb';
// import {GetInfoRequest, GetInfoResponse} from '../generated/flipdot_pb';
import {FlipAppsClient} from '../generated/flipapps_pb_service';
import {grpc} from '@improbable-eng/grpc-web';

export class Client {
    // Token used to send authenticated gRPC messages with
    private token: string | null = null;

    // gRPC client
    private client: FlipAppsClient;

    // Error returned by server
    private err: any = null;

    constructor(domain: string) {
        // Create a flipapps client
        this.client = new FlipAppsClient(domain);
    }

    public authenticate(password: string, callback: (response: any) => void) {
        // Construct a request
        const request = new AuthenticateRequest();
        request.setPassword(password);
        // Send the request
        this.client.authenticate(request, new grpc.Metadata(), (err, response) => {
            this.handle(err, response, (r) => {
                // Save the token
                this.token = r.getToken();
                // Execute user callback
                callback(r);
            });
        });
    }

    public sendTextMessage(from: string, text: string, callback: (response: any) => void) {
        console.assert(this.isAuthenticated);
        // Construct a request
        const request = new MessageRequest();
        request.setFrom(from);
        request.setText(text);
        // Send the request
        this.client.sendMessage(request, new grpc.Metadata({token: this.token as string}), (err, response) => {
            // Handle any errors, or execute user callback
            this.handle(err, response, callback);
        });
    }

    get error(): any {
        return this.err;
    }

    get isAuthenticated(): boolean {
        return this.token != null;
    }

    private handle(err: any, response: any, callback: (response: any) => void) {
        // First ensure we've not errored out
        if (err != null) {
            this.err = err;
            this.handleError();
            return;
        } else {
            // Reset error
            this.err = null;
        }
        // Otherwise continue
        callback(response);
    }

    private handleError() {
        // Reset the token in the event of an authentication failure
        if (this.err.code === grpc.Code.Unauthenticated) {
            this.token = null;
        }
    }
}
