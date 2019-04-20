import { mount, Wrapper } from '@vue/test-utils';
import { expect } from 'chai';
import Message from '@/components/Message.vue';
import { Client } from '@/ts/client';
import { Interpreter } from 'xstate/lib/interpreter';
import { mock, instance, when, anything, verify } from 'ts-mockito';
import { CombinedVueInstance, Vue } from 'vue/types/vue';
import { grpc } from '@improbable-eng/grpc-web';
import { StateSchema, EventObject } from 'xstate';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';

describe('Message.vue', () => {
    let mockedClient: Client;
    let client: Client;

    let mockedFsm: Interpreter<{}, StateSchema, EventObject>;
    let fsm: Interpreter<{}, StateSchema, EventObject>;

    let wrapper: Wrapper<CombinedVueInstance<Message, object, object, object, Record<never, any>>>;

    beforeEach(() => {
        // create an extended `Vue` constructor
        const localVue = createLocalVue();
        // install plugins as normal
        localVue.use(BootstrapVue);
        // Create a mock client
        mockedClient = mock(Client);
        client = instance(mockedClient);
        // Create a mock fsm
        mockedFsm = mock(Interpreter);
        fsm = instance(mockedFsm);
        // Create object under test
        wrapper = mount(Message, {
            localVue,
            propsData: {
                client,
                fsm,
            },
        });
    });

    it('sends successfully', () => {
        const sender = 'briggySmalls';
        const text = 'it was all a dream';
        // Configure client
        when(mockedClient.error).thenReturn(null);
        when(mockedClient.sendTextMessage(sender, text, anything())).thenCall(
            (frm: string, msg: string, callback: (response: any) => void) => {
                // Execute callback
                callback(null);
            },
        );
        // Send message
        wrapper.find('form #sender-field').setValue(sender);
        wrapper.find('form #text-field').setValue(text);
        wrapper.find('form #message-submit').trigger('submit');
        // Set a password
        verify(mockedClient.sendTextMessage(sender, text, anything())).once();
        verify(mockedFsm.send('SENT')).once();
    });

    it('reauthenticates on expiry', () => {
        const sender = 'briggySmalls';
        const text = 'it was all a dream';
        // Configure client
        when(mockedClient.error).thenReturn(null);
        when(mockedClient.sendTextMessage(sender, text, anything())).thenCall(
            (sndr: string, txt: string, callback: (response: any) => void) => {
                // Indicate an error occurred
                when(mockedClient.error).thenReturn({
                    code: grpc.Code.Unauthenticated,
                    message: 'Failed to login',
                });
                // Execute callback
                callback(null);
            },
        );
        // Send message
        wrapper.find('form #sender-field').setValue(sender);
        wrapper.find('form #text-field').setValue(text);
        wrapper.find('form #message-submit').trigger('submit');
        // Assert state machine
        verify(mockedFsm.send('REAUTH')).once();
    });

    it('proceeds on error', () => {
        const sender = 'briggySmalls';
        const text = 'it was all a dream';
        // Configure client
        when(mockedClient.error).thenReturn(null);
        when(mockedClient.sendTextMessage(sender, text, anything())).thenCall(
            (sndr: string, txt: string, callback: (response: any) => void) => {
                // Indicate an error occurred
                when(mockedClient.error).thenReturn({
                    code: grpc.Code.Unavailable,
                    message: 'Failed to find server',
                });
                // Execute callback
                callback(null);
            },
        );
        // Send message
        wrapper.find('form #sender-field').setValue(sender);
        wrapper.find('form #text-field').setValue(text);
        wrapper.find('form #message-submit').trigger('submit');
        // Assert state machine
        verify(mockedFsm.send('SENT')).once();
    });
});
