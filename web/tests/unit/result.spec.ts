import { mount, Wrapper } from '@vue/test-utils';
import { expect } from 'chai';
import Result from '@/components/Result.vue';
import { Client } from '@/ts/client';
import { Interpreter } from 'xstate/lib/interpreter';
import { mock, instance, when, anything, verify } from 'ts-mockito';
import { CombinedVueInstance, Vue } from 'vue/types/vue';
import { grpc } from '@improbable-eng/grpc-web';
import { StateSchema, EventObject } from 'xstate';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';

describe('Result.vue', () => {
    let mockedClient: Client;
    let client: Client;

    let mockedFsm: Interpreter<{}, StateSchema, EventObject>;
    let fsm: Interpreter<{}, StateSchema, EventObject>;

    let wrapper: Wrapper<CombinedVueInstance<Result, object, object, object, Record<never, any>>>;

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
        wrapper = mount(Result, {
            localVue,
            propsData: {
                client,
                fsm,
            },
        });
    });

    it('displays sent and resend', () => {
        // Configure client to indicate no error
        when(mockedClient.error).thenReturn(null);
        // Check success is displayed
        wrapper.find('.alert').classes().includes('alert-success');
        // Indicate we wish to send another message
        wrapper.find('button').trigger('click');
        verify(mockedFsm.send('NEW')).once();
    });

    it('displays error and resend', () => {
        // Configure client to indicate no error
        when(mockedClient.error).thenReturn({
            code: grpc.Code.Unavailable,
            message: 'Couldn\'t find server',
        });
        // Check error is displayed
        wrapper.find('.alert').classes().includes('alert-danger');
        // Indicate we wish to send another message
        wrapper.find('button').trigger('click');
        verify(mockedFsm.send('NEW')).once();
    });
});
